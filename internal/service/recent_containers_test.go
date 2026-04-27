package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRecentContainers_Record 测试基本插入次序：最新插入的应在最前面
func TestRecentContainers_Record(t *testing.T) {
	r := NewRecentContainers()

	r.Record("a")
	r.Record("b")
	r.Record("c")

	ids := r.GetIDs()
	assert.Equal(t, []string{"c", "b", "a"}, ids)
}

// TestRecentContainers_RecordPromotes 测试重复插入：已存在的元素会被提升到头部
func TestRecentContainers_RecordPromotes(t *testing.T) {
	r := NewRecentContainers()

	r.Record("a")
	r.Record("b")
	r.Record("c")
	r.Record("a") // a 应被提升到最前面

	ids := r.GetIDs()
	assert.Equal(t, []string{"a", "c", "b"}, ids)
}

// TestRecentContainers_MaxSize 测试容量上限：超过 5 个时淘汰最久未使用的元素
func TestRecentContainers_MaxSize(t *testing.T) {
	r := NewRecentContainers()

	for i := 0; i < 7; i++ {
		r.Record(string(rune('a' + i)))
	}

	ids := r.GetIDs()
	assert.Len(t, ids, maxRecent) // 最多保留 5 个
	assert.Equal(t, "g", ids[0])
	assert.Equal(t, "f", ids[1])
	assert.Equal(t, "e", ids[2])
	assert.NotContains(t, ids, "a") // a 和 b 已被淘汰
	assert.NotContains(t, ids, "b")
}

// TestRecentContainers_Remove 测试删除元素
func TestRecentContainers_Remove(t *testing.T) {
	r := NewRecentContainers()

	r.Record("a")
	r.Record("b")
	r.Record("c")
	r.Remove("b")

	ids := r.GetIDs()
	assert.Equal(t, []string{"c", "a"}, ids)
	assert.Len(t, ids, 2)
}

// TestRecentContainers_RemoveNonExistent 测试删除不存在的元素：不应对列表产生影响
func TestRecentContainers_RemoveNonExistent(t *testing.T) {
	r := NewRecentContainers()

	r.Record("a")
	r.Record("b")
	r.Remove("z") // 删除不存在的元素

	ids := r.GetIDs()
	assert.Equal(t, []string{"b", "a"}, ids)
}

// TestRecentContainers_Empty 测试空列表
func TestRecentContainers_Empty(t *testing.T) {
	r := NewRecentContainers()
	ids := r.GetIDs()
	assert.Empty(t, ids)
}
