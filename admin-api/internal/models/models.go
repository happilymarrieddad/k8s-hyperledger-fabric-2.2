package models

import (
	"admin-api/internal/k8client"
	"sync"

	"xorm.io/xorm"
)

type GlobalModels interface {
	Organizations() Organizations
}

func NewGlobalModels(db *xorm.Engine, k8c k8client.Client) GlobalModels {
	return &globalModels{db, k8c, make(map[string]interface{}), &sync.Mutex{}}
}

type globalModels struct {
	db     *xorm.Engine
	k8c    k8client.Client
	models map[string]interface{}
	mutex  *sync.Mutex
}

func (g *globalModels) getModel(key string, gen func(*xorm.Engine, k8client.Client) interface{}) interface{} {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	val, exists := g.models[key]
	if exists {
		return val
	}

	g.models[key] = gen(g.db, g.k8c)
	return g.models[key]
}

func (g *globalModels) Organizations() Organizations {
	return g.getModel("organizations", func(db *xorm.Engine, k8c k8client.Client) interface{} { return NewOrganizations(db, k8c) }).(Organizations)
}
