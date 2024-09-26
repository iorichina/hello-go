package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("两数相加：")
	{
		l1 := &ListNode{2, &ListNode{4, &ListNode{3, nil}}}
		l2 := &ListNode{5, &ListNode{6, &ListNode{4, nil}}}
		v := addTwoNumbers(l1, l2)
		fmt.Println(strListNode(v))
	}
	{
		l1 := &ListNode{0, nil}
		l2 := &ListNode{0, nil}
		fmt.Println(strListNode(addTwoNumbers(l1, l2)))
	}
	{
		l1 := &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, nil}}}}}}}}}
		l2 := &ListNode{9, &ListNode{9, &ListNode{9, &ListNode{9, nil}}}}
		fmt.Println(strListNode(addTwoNumbers(l1, l2)))
	}
}

func strListNode(l *ListNode) string {
	sj := []string{}
	for nil != l {
		sj = append(sj, strconv.Itoa((l.Val)))
		l = l.Next
	}
	return "[" + strings.Join(sj, ",") + "]"
}

type ListNode struct {
	Val  int
	Next *ListNode
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	head, tail := l1, l1
	t := 0
	for nil != l1 && nil != l2 {
		t += l1.Val + l2.Val
		if t > 9 {
			tail.Val = t - 10
			t = 1
		} else {
			tail.Val = t
			t = 0
		}
		l1 = l1.Next
		l2 = l2.Next
		if nil != l1 {
			tail.Next = l1
			tail = tail.Next
		} else if nil != l2 {
			tail.Next = l2
			tail = tail.Next
		}
	}
	if nil != l1 {
		for nil != l1 {
			t += l1.Val
			if t > 9 {
				tail.Val = t - 10
				t = 1
			} else {
				tail.Val = t
				t = 0
			}
			l1 = l1.Next
			if nil != l1 {
				tail.Next = l1
				tail = tail.Next
			}
		}
	}
	if nil != l2 {
		for nil != l2 {
			t += l2.Val
			if t > 9 {
				tail.Val = t - 10
				t = 1
			} else {
				tail.Val = t
				t = 0
			}
			l2 = l2.Next
			if nil != l2 {
				tail.Next = l2
				tail = tail.Next
			}
		}
	}
	if t > 0 {
		// tail.Next = new(ListNode)
		tail.Next = &ListNode{}
		tail.Next.Val = t
	}
	return head
}
