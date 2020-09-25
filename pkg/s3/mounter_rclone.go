package s3

import (
	"fmt"
	"os"
)

// Implements Mounter
type rcloneMounter struct {
	bucket          *bucket
	url             string
	region          string
	accessKeyID     string
	secretAccessKey string
	args            []string
}

const (
	rcloneCmd = "rclone"
)

func newRcloneMounter(b *bucket, cfg *Config) (Mounter, error) {
	return &rcloneMounter{
		bucket:          b,
		url:             cfg.Endpoint,
		region:          cfg.Region,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
		args:            b.MounterArgs,
	}, nil
}

func (rclone *rcloneMounter) Stage(stageTarget string) error {
	return nil
}

func (rclone *rcloneMounter) Unstage(stageTarget string) error {
	return nil
}

func (rclone *rcloneMounter) Mount(source string, target string) error {
	args := []string{
		"mount",
		fmt.Sprintf(":s3:%s/%s", rclone.bucket.Name, rclone.bucket.FSPath),
		fmt.Sprintf("%s", target),
		"--daemon",
		"--s3-provider=AWS",
		"--s3-env-auth=true",
		fmt.Sprintf("--s3-region=%s", rclone.region),
		fmt.Sprintf("--s3-endpoint=%s", rclone.url),
		"--allow-other",
	}
	args = append(args, rclone.args...)
	os.Setenv("AWS_ACCESS_KEY_ID", rclone.accessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", rclone.secretAccessKey)
	return fuseMount(target, rcloneCmd, args)
}
