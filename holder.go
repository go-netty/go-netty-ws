/*
 * Copyright 2023 the go-netty project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nettyws

import (
	"fmt"
	"sync"

	"github.com/go-netty/go-netty"
)

// newChannelHolder create a new ChannelHolder with initial capacity
func newChannelHolder(capacity int) netty.ChannelHolder {
	return &channelHolder{channels: make(map[int64]netty.Channel, capacity)}
}

type channelHolder struct {
	channels map[int64]netty.Channel
	mutex    sync.Mutex
}

func (c *channelHolder) HandleActive(ctx netty.ActiveContext) {
	c.addChannel(ctx.Channel())
	ctx.HandleActive()
}

func (c *channelHolder) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	c.delChannel(ctx.Channel())
	ctx.HandleInactive(ex)
}

func (c *channelHolder) CloseAll(err error) {
	c.mutex.Lock()
	channels := c.channels
	c.channels = make(map[int64]netty.Channel, 1024)
	c.mutex.Unlock()

	// close reason
	wse, ok := err.(ClosedError)

	for _, ch := range channels {
		if ok {
			_ = ch.Transport().(wsc).WriteClose(wse.Code, wse.Reason)
		}
		ch.Close(err)
	}
}

func (c *channelHolder) addChannel(ch netty.Channel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	id := ch.ID()
	if _, ok := c.channels[id]; ok {
		panic(fmt.Errorf("duplicate channel: %d", id))
	}
	c.channels[id] = ch
}

func (c *channelHolder) delChannel(ch netty.Channel) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.channels, ch.ID())
}
