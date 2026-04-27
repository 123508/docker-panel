package service

import (
	"container/list"
	"sync"
)

// 最近操作容器列表的最大容量
const maxRecent = 5

// RecentContainers 线程安全的 LRU 结构，用于存储最近操作过的容器 ID。
// 内部使用双向链表 + map 实现，Record/Remove/GetIDs 均为 O(1) 时间复杂度。
type RecentContainers struct {
	mu    sync.Mutex
	order *list.List               // 双向链表，头部为最新，尾部为最旧
	index map[string]*list.Element // 从容器 ID 到链表节点的映射
}

// NewRecentContainers 创建并返回一个新的 RecentContainers 实例。
func NewRecentContainers() *RecentContainers {
	return &RecentContainers{
		order: list.New(),
		index: make(map[string]*list.Element),
	}
}

// Record 将指定容器 ID 标记为最近操作。
// 若该 ID 已存在则将其移到链表头部；若链表长度超过 maxRecent，则淘汰尾部（最久未使用）的元素。
func (r *RecentContainers) Record(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 已存在则提升到头部
	if el, ok := r.index[id]; ok {
		r.order.MoveToFront(el)
		return
	}

	// 新元素插入链表头部
	el := r.order.PushFront(id)
	r.index[id] = el

	// 超过容量上限时淘汰尾部
	for r.order.Len() > maxRecent {
		old := r.order.Back()
		if old != nil {
			r.order.Remove(old)
			delete(r.index, old.Value.(string))
		}
	}
}

// Remove 从最近操作列表中删除指定容器 ID。
// 主要用于容器被删除时清理记录，避免展示已不存在的容器。
func (r *RecentContainers) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if el, ok := r.index[id]; ok {
		r.order.Remove(el)
		delete(r.index, id)
	}
}

// GetIDs 返回所有容器 ID 的副本，按访问时间从近到远排序。
func (r *RecentContainers) GetIDs() []string {
	r.mu.Lock()
	defer r.mu.Unlock()

	ids := make([]string, 0, r.order.Len())
	for el := r.order.Front(); el != nil; el = el.Next() {
		ids = append(ids, el.Value.(string))
	}
	return ids
}
