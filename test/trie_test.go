package test

import (
	"trie/trie"
	"log"
	"testing"
)

func TestMain(t *testing.T) {
	TrieTree := trie.NewTrie()
	list := []string{"关键词","过滤"}
	for _, item := range list {
		TrieTree.Insert(item)
	}

	// 查找关键词
	isExist, l := TrieTree.Search("我是关键词，你想过滤我吗？")
	log.Println("是否存在关键词:",isExist,"关键词个数:",len(l))
	for _, v := range l {
		log.Println("关键词:",string(v))
	}

	// 替换关键词 可直接使用
	new := TrieTree.Replace("我是关键词，你想过滤我吗？")
	log.Println(new)
	
}