data "template_file" "cloud_config" {
  template = "${file("${path.module}/cloud_config.tpl")}"
  vars = {
    username = "luther"
  }
}

data "template_file" "meta_data" {
  template = "${file("${path.module}/meta_data.tpl")}"
  vars = {
    hostname = "test-terraform"
  }
}

data "template_cloudinit_config" "user_data" {
  gzip          = false
  base64_encode = false

  # Main cloud-config configuration file.
  part {
    filename     = "init.cfg"
    content_type = "text/cloud-config"
    content      = "${data.template_file.cloud_config.rendered}"
  }

  part {
    content_type = "text/x-shellscript"
    content      = "ls -la"
  }

  part {
    content_type = "text/x-shellscript"
    content      = "cd /tmp && ls -la"
  }
}

resource "local_file" "user_data" {
    content  = data.template_cloudinit_config.user_data.rendered
    filename = "${path.module}/user-data"
}

resource "local_file" "meta_data" {
    content  = data.template_file.meta_data.rendered
    filename = "${path.module}/meta-data"
}

resource "null_resource" "echo" {
  provisioner "local-exec" {
    command = "sleep 2 && quiso && vagrant up"
  }
}
