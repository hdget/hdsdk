// Package parallel
// 实现了可优雅关闭退出所有的actors
// Group
//   | ---> 执行actor
//   | ---> 退出处理
// 其类似errgroup, 但其不要求actor goroutines去理解上下文的意思，这样它可以在更多的场景下使用.
// 例如, 处理net.Listeners的connections, 或者从关闭的io.Reader中读取输入
package parallel

type actor struct {
	execute   func() error
	interrupt func(error)
}

// Group 收集一系列actor并且同时运行它们
// 当其中某一个actor执行完，其他actors都会被中断执行
type Group struct {
	actors []actor
}

// Add 添加actor到group中。
// 每个actor都必须被中断函数是抢占控制，换句话说，当中断函数执行，actor执行函数必须退出
// 同时，必须保证执行函数退出后再调用中断函数也是安全的
//
// 首个执行完成的actor会中断所有正在运行的actor, 通过将error传给中断函数，最后被Run返回
func (g *Group) Add(execute func() error, interrupt func(error)) {
	g.actors = append(g.actors, actor{execute, interrupt})
}

// Run 并行运行所有actors
// 当首个actor执行完成，其他的actor都会被中断
// 只有所有actors退出后Run才返回，否则会一直会阻塞
// Run最后返回时会返回首个退出的actor返回的error
func (g *Group) Run() error {
	if len(g.actors) == 0 {
		return nil
	}

	// Run each actor.
	errors := make(chan error, len(g.actors))
	for _, a := range g.actors {
		go func(a actor) {
			errors <- a.execute()
		}(a)
	}

	// Wait for the first actor to stop.
	err := <-errors

	// Signal all actors to stop.
	for _, a := range g.actors {
		a.interrupt(err)
	}

	// Wait for all actors to stop.
	for i := 1; i < cap(errors); i++ {
		<-errors
	}

	// Return the original error.
	return err
}
