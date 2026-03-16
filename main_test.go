package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListenAddr_ValidPort(t *testing.T) {
	addr, err := listenAddr(8500)
	require.NoError(t, err)
	require.Equal(t, ":8500", addr)
}

func TestListenAddr_InvalidPort(t *testing.T) {
	_, err := listenAddr(0)
	require.Error(t, err)
}
