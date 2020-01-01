# quiso, the quick ISO builder for [cloud-init](https://cloudinit.readthedocs.io/en/latest/)
Simple cross-platform CLI with cloud-init with presets designed to make building images painless and simple, 
designed to work everywhere for quick iteration on cloud-init development. There are also unsupported clouds where the 
[NoCloud](https://cloudinit.readthedocs.io/en/latest/topics/datasources/nocloud.html) method of injecting user-data
is the best method of using cloud-init.

This CLI was inspired by work done for 
[rancher-machine](https://github.com/rancher/machine/blob/master/drivers/vmwarevsphere/cloudinit.go#L127-L201) to 
accommodate vSphere without relying on the [guestinfo](https://github.com/vmware/cloud-init-vmware-guestinfo) 
implementation which requires VMWare Tools. Other examples of this same solution can be found in Canonical's 
[multipass](https://github.com/canonical/multipass/blob/master/src/iso/cloud_init_iso.cpp) and countless 
[bash](https://github.com/frederickding/Cloud-Init-ISO) implementations which rely on 
[genisoimage](https://www.alextomkins.com/2016/09/testing-cloud-init-with-vagrant/),
[mkisofs](https://januz.nl/mac/generate-coreos-configdrive-iso-based-cloud-config-file) or 
[hdiutil](https://serverfault.com/questions/576615/os-x-equivalents-for-ubuntus-genisoimage-and-qemu-img). This 
tool takes these ideas and coalesces them into a single tool with easy defaults allowing any devops engineer in any
environment (mac, win, linux) to call `quiso` and immediately get an ISO.

## Releases
Find everything you need to start in your architecture on the 
[releases](https://github.com/luthermonson/quiso/releases) page.

## Usage
`quiso` was designed with simple defaults, run in any directory containing a `user-data` and `meta-data` file and 
it will automatically create `user-data.iso`.   

```shell
$ tree
   .
   ├── meta-data
   └── user-data
$ quiso
$ tree
   .
   ├── meta-data
   └── user-data
   └── user-data.iso
```

Calling `quiso` with no parameters is equivalent of 

```
quiso build --user-data ./user-data --meta-data ./meta-data --output ./user-data.iso 
``` 

## Examples
Check out the [./examples](https://github.com/luthermonson/quiso/tree/master/examples) directory for how to 
use [vagrant](https://www.vagrantup.com/downloads.html) and 
[virtualbox](https://www.virtualbox.org/wiki/Downloads) to test ISOs locally and how to use `quiso` with
terraform's [template_cloudinit_config](https://www.terraform.io/docs/providers/template/d/cloudinit_config.html).

## todo
Looking for suggestions, leave an [issue](https://github.com/luthermonson/quiso/issues) with a request or fork and 
submit a pull request.