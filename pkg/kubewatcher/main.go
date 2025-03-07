/*
Copyright 2016 The Fission Authors.

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

package kubewatcher

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/fission/fission/pkg/crd"
	"github.com/fission/fission/pkg/publisher"
)

func Start(ctx context.Context, logger *zap.Logger, routerUrl string) error {
	clientGen := crd.NewClientGenerator()
	fissionClient, err := clientGen.GetFissionClient()
	if err != nil {
		return errors.Wrap(err, "failed to get fission client")
	}
	kubeClient, err := clientGen.GetKubernetesClient()
	if err != nil {
		return errors.Wrap(err, "failed to get kubernetes client")
	}

	err = crd.WaitForCRDs(ctx, logger, fissionClient)
	if err != nil {
		return errors.Wrap(err, "error waiting for CRDs")
	}

	poster := publisher.MakeWebhookPublisher(logger, routerUrl)
	kubeWatch := MakeKubeWatcher(ctx, logger, kubeClient, poster)
	ws := MakeWatchSync(ctx, logger, fissionClient, kubeWatch)
	ws.Run(ctx)

	return nil
}
