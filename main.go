package main

import (
	"net"

	"github.com/harborapp/harbor-client/gradlew"
	"github.com/harborapp/harbor-client/manifest"
	"github.com/harborapp/harbor-client/project"
	log "github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type builderServer struct{}

func (bsrv *builderServer) Build(ctx context.Context, job *BuildJobRequest) (*BuildJobResponse, error) {
	builder, err := gradlew.New(
		gradlew.WithPath(job.Gradlew),
		gradlew.WithProjPath(job.ProjPath),
		gradlew.WithOutputPath(job.Output),
	)

	if err != nil {
		return nil, err
	}

	pkgr := manifest.New(
		manifest.WithPath(job.Manifest),
	)

	p, err := project.New(builder, pkgr)
	if err != nil {
		return nil, err
	}

	data, err := p.BuildProject(job.Task)
	if err != nil {
		return nil, err
	}

	var apks []*Apk
	for _, v := range data {
		apks = append(apks, &Apk{
			Path:    v.Path,
			RawSize: v.RawSize,
			Size:    v.Size,
		})
	}

	return &BuildJobResponse{
		Apks:    apks,
		Success: true,
	}, nil
}

func main() {
	srv := grpc.NewServer()
	bsrv := &builderServer{}
	RegisterBuilderServer(srv, bsrv)

	// TODO: Remove hardcoding.
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("Could not listen on port")
	}

	log.Fatal(srv.Serve(lis))
}
