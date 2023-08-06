package leader_election

import (
	"context"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type EtcdLeaderElection struct {
	LeaderElectionConfig
	session  *concurrency.Session
	election *concurrency.Election
}

type LeaderElectionConfig struct {
	// etcdClient 客户端
	EtcdClient *clientv3.Client
	// leaseSeconds 租约时间：当如果过了LeaseSeconds时间后，会强制重新选举
	LeaseSeconds int
	// Callbacks 选主不同阶段的回调方法
	Callbacks LeaderCallbacks
	// ElectionName 选主实例名
	ElectionName string
	// Identity 唯一标示，选主机制底层是分布式锁
	Identity string
}

// LeaderCallbacks 选主机制的回调方法
type LeaderCallbacks struct {
	// OnStartedLeading is called when a LeaderElector client starts leading
	// OnStartedLeading 当获取到锁时，执行的回调
	OnStartedLeading func(context.Context)
	// OnStoppedLeading is called when a LeaderElector client stops leading
	// OnStoppedLeading 当释放锁时，执行的回调
	OnStoppedLeading func()
	// OnNewLeader 当重新选主后，执行的回调
	OnNewLeader func(identity string)
}

func NewLeaderElection(config LeaderElectionConfig) (*EtcdLeaderElection, error) {
	// 生成session
	session, err := concurrency.NewSession(config.EtcdClient, concurrency.WithTTL(config.LeaseSeconds))
	if err != nil {
		return nil, err
	}
	// 生成选举实例
	election := concurrency.NewElection(session, config.ElectionName)

	return &EtcdLeaderElection{
		LeaderElectionConfig: config,
		session:              session,
		election:             election,
	}, nil
}

// Run 执行选主
func (le *EtcdLeaderElection) Run(ctx context.Context) error {
	// 退出Run方法，即释放锁
	defer func() {
		if le.Callbacks.OnStoppedLeading != nil {
			le.Callbacks.OnStoppedLeading()
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 异步观察选主是否变更
	go le.observe(ctx)

	// 启动选举，不断去获取锁，如果没有获取到锁，会阻塞
	if err := le.election.Campaign(ctx, le.Identity); err != nil {
		return err
	}
	// 当获取到锁后，执行回调方法
	if le.Callbacks.OnStartedLeading != nil {
		le.Callbacks.OnStartedLeading(ctx)
	}

	return nil
}

// observe 观察是否变更选主
func (le *EtcdLeaderElection) observe(ctx context.Context) {
	if le.Callbacks.OnNewLeader == nil {
		return
	}

	ch := le.election.Observe(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case resp, ok := <-ch:
			if !ok {
				return
			}

			if len(resp.Kvs) == 0 {
				continue
			}

			leader := string(resp.Kvs[0].Value)
			if leader != le.Identity {
				go le.Callbacks.OnNewLeader(leader)
			}
		}
	}
}

func (le *EtcdLeaderElection) Close() error {
	return le.session.Close()
}


