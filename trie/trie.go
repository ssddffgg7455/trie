package trie

import (
	"container/list"
	"strings"
	"sync"
)

type Trie struct {
	next  map[rune]*Trie
	isEnd bool
}

func NewTrie() (trie Trie) {
	trie = Trie{}
	trie.next = make(map[rune]*Trie)
	return
}

// 未加锁 启动一次性加载
func (this *Trie) Insert(word string) {
	node := this
	for _, v := range word {

		if node.next[v] == nil {
			trie := new(Trie)
			trie.next = make(map[rune]*Trie)
			node.next[v] = trie
		}
		node = node.next[v]
	}
	node.isEnd = true
}

type Replaced struct {
	Lock sync.Mutex
	List [][]rune
}

func (this *Trie) Search(words string) (isExist bool, replacedList [][]rune) {
	isExist = false
	l := list.New()

	for _, c := range words {
		l.PushBack(c)
	}
	// 并发查询
	var wg sync.WaitGroup
	var rep = new(Replaced)
	for sub := l.Front(); sub != nil; sub = sub.Next() {
		_, ok := sub.Value.(rune)
		if ok {
			wg.Add(1)
			go this._search(&wg, sub, rep)
		}
	}
	wg.Wait()

	if len(rep.List) > 0 {
		isExist = true
		replacedList = rep.List
	}
	return
}

func (this *Trie) _search(wg *sync.WaitGroup, sub *list.Element, rep *Replaced) bool {
	node := this
	oldWords := []rune{} 

	for {
		if node.next[sub.Value.(rune)] == nil {
			break
		} else {
			node = node.next[sub.Value.(rune)]
		}

		oldWords = append(oldWords, sub.Value.(rune))
		
		if sub = sub.Next(); sub == nil {
			break
		}
	}	
	
	rep.Lock.Lock()
	defer rep.Lock.Unlock()
	if len(oldWords) > 0 && node.isEnd {
		rep.List = append(rep.List, oldWords)
	}
	
	wg.Done()
	return node.isEnd
}

func (this *Trie) Replace(words string) string {
	isExist, list := this.Search(words)

	if isExist {
		for _, oldWords := range list {
			newWords := strings.Repeat("*", len(oldWords))
			words = strings.Replace(words, string(oldWords), newWords, -1)
		}
	}
	
	return words
}
