module github.com/forensicanalysis/fslib

go 1.12

require (
	github.com/forensicanalysis/go-vss v1.2.0
	github.com/h2non/filetype v1.0.12
	github.com/ledongthuc/pdf v0.0.0-20200323191019-23c5852adbd2
	github.com/spf13/afero v1.2.3-0.20200410222221-ceb6a5e37254
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/xlab/treeprint v1.0.0
	golang.org/x/sys v0.0.0-20200413165638-669c56c373c4
	gopkg.in/djherbis/times.v1 v1.2.0
	www.velocidex.com/golang/go-ntfs v0.0.0-20200530234845-9c557c0b9eec
)

replace github.com/forensicanalysis/go-vss => ../go-vss
