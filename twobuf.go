// Copyright 2020 guonaihong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httparser

// 设计思路:

// TwoBuf是为了减少内存拷贝设计的一种数据结构
// 适用的场景有httparser解析器，
// httparser解析器返回值有两种情况
// 1.有时等于输入数据
// 2.有时不等于输入数据，这时候标志还有一些数据下次要重新送入
// 怎么处理这种需要拷贝数据的场景？正常想法是把溢出数据再和新数据append一下组成新的buf送入解析器
//
// TwoBuf相比append的方式是怎么处理溢出情况?
// 1.首先:twobuf的内存布局是->实际空间是申请空间的两倍
//	   TwoBuf由两部分组成, buf = left + right
//	   right:占buf一半空间存放实际数据
//	   left: 占buf一半部分存放溢出数据
// 2.然后:
// 只要把上次溢出的数据拷贝至左边，这次要读入的数据放至右边。
// 记录左边的offset直接buf[left:] 就是新的buffer数据

// 为什么这么设计:
// 溢出(未解析的数据)是概率比较少见，而且溢出的数据不会太长，
// 重新分配大块内存分配数据有点浪费

type TwoBuf struct {
	buf  []byte
	mid  int // 中间位置
	left int
}

// 新建twobuf, size表明单个块的大小
func NewTwoBuf(size int) *TwoBuf {
	return &TwoBuf{buf: make([]byte, size*2), mid: size, left: size}
}

// 获取右边buf
func (t *TwoBuf) Right() []byte {
	return t.buf[t.mid:]
}

// 获取全部buf, 如果没有溢出数据，就取右边
// 如果有溢出数据，就把溢出数据也包含进来
func (t *TwoBuf) All(right int) []byte {
	return t.buf[t.left : t.mid+right]
}

// 移动到左边buf
// 如果写入数据超过left空间，会直接panic
func (t *TwoBuf) MoveLeft(leftBuf []byte) {
	t.left = t.mid - len(leftBuf)
	if t.left < 0 {
		panic("TwoBuf:The copied buf is too large")
	}

	copy(t.buf[t.left:], leftBuf)
}

func (t *TwoBuf) Reset() {
	t.left = t.mid
}
