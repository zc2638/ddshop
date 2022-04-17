// Package music
// Created by zc on 2022/4/17.
package music

import (
	"testing"

	"github.com/zc2638/ddshop/asserts"
)

func TestNewMP3(t *testing.T) {
	player, err := NewMP3(asserts.NoticeMP3, 180)
	if err != nil {
		return
	}
	if err := player.Play(); err != nil {
		t.Fatal(err)
	}
}
