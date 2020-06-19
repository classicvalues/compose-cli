/*
   Copyright 2020 Docker, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package proxy

import (
	"context"

	"github.com/docker/api/config"
	"github.com/docker/api/context/store"
	contextsv1 "github.com/docker/api/protos/contexts/v1"
)

type contextsProxy struct {
	configDir string
}

func (cp *contextsProxy) SetCurrent(ctx context.Context, request *contextsv1.SetCurrentRequest) (*contextsv1.SetCurrentResponse, error) {
	if err := config.WriteCurrentContext(cp.configDir, request.GetName()); err != nil {
		return &contextsv1.SetCurrentResponse{}, err
	}

	return &contextsv1.SetCurrentResponse{}, nil
}

func (cp *contextsProxy) List(ctx context.Context, request *contextsv1.ListRequest) (*contextsv1.ListResponse, error) {
	s := store.ContextStore(ctx)
	configFile, err := config.LoadFile(cp.configDir)
	if err != nil {
		return nil, err
	}
	contexts, err := s.List()
	if err != nil {
		return &contextsv1.ListResponse{}, err
	}

	result := &contextsv1.ListResponse{}

	for _, c := range contexts {
		result.Contexts = append(result.Contexts, &contextsv1.Context{
			Name:        c.Name,
			ContextType: c.Type(),
			Current:     c.Name == configFile.CurrentContext,
		})
	}

	return result, nil
}
