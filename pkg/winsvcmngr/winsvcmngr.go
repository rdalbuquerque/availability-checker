package winsvcmngr

import (
	"golang.org/x/sys/windows/svc/mgr"
)

type WinSvcMngr interface {
	Connect() error
	Disconnect() error
	OpenService(name string) (WinSvc, error)
}

type DefaultWinSvcMngr struct {
	mngr *mgr.Mgr
}

func (svcmngr *DefaultWinSvcMngr) Connect() error {
	mngr, err := mgr.Connect()
	if err != nil {
		return err
	}
	svcmngr.mngr = mngr
	return nil
}

func (svcmngr *DefaultWinSvcMngr) Disconnect() error {
	return svcmngr.mngr.Disconnect()
}

func (svcmngr *DefaultWinSvcMngr) OpenService(name string) (WinSvc, error) {
	svc, err := svcmngr.mngr.OpenService(name)
	if err != nil {
		return nil, err
	}
	return &DefaultWinSvc{svc: svc}, err
}

type WinSvc interface {
	Close() error
	Start() error
}

type DefaultWinSvc struct {
	svc *mgr.Service
}

func (svc *DefaultWinSvc) Start() error {
	return svc.svc.Start()
}

func (svc *DefaultWinSvc) Close() error {
	return svc.svc.Close()
}
