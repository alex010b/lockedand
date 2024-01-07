package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
)

func Listener(events chan VPS, conn *websocket.Conn) {
	for {
		select {
		case event := <-events:

			SendWs(conn, "configuring ip and port")
			fmt.Println("NIGGERS", event)
			currentport, _ := os.ReadFile("currentport")
			fmt.Println(string(currentport), event)
			currentporti, _ := strconv.Atoi(string(currentport))
			Currentport := fmt.Sprint(currentporti)
			currentporti++
			os.WriteFile("currentport", []byte(strconv.Itoa(currentporti)), 0644)
			currentip, _ := os.ReadFile("currentip")
			fmt.Println(string(currentip), event)
			currentipi, _ := strconv.Atoi(string(currentip))
			Currentip := fmt.Sprint(currentipi)
			currentporti++
			os.WriteFile("currentip", []byte(strconv.Itoa(currentipi)), 0644)

			SendWs(conn, "generating installation config")

			preseed := CreatePreseed(event.Password, event.Hostname)

			preseedPath := "database/" + event.Username + "/vm/preseed.cfg"
			if err := os.WriteFile(preseedPath, []byte(preseed), 0644); err != nil {
				fmt.Println("couldnt write preseed file", err)
			}

			xmlPath := fmt.Sprintf("database/%s/vm/%s.xml", event.Username, event.Os)
			diskPath := fmt.Sprintf("database/%s/vm/%s.qcow2", event.Username, event.Os)
			isoPath := fmt.Sprintf("database/%s/vm/%s_iso", event.Username, event.Os)
			customIsoPath := fmt.Sprintf("database/%s/vm/custom_%s_iso", event.Username, event.Os)
			xorrisoIsoPath := fmt.Sprintf("database/%s/vm/custom_%s.iso", event.Username, event.Os)

			SendWs(conn, "creating qcow2 disk")

			command1 := fmt.Sprintf(`mkdir %s; sudo mount -o loop %s.iso %s; mkdir %s; cp -r %s/* %s/; cp %s %s/; sudo umount %s; xorriso -as mkisofs -r -J -l -isohybrid-mbr /usr/lib/ISOLINUX/isohdpfx.bin -c isolinux/boot.cat -b isolinux/isolinux.bin -no-emul-boot -boot-load-size 4 -boot-info-table -eltorito-alt-boot -e boot/grub/efi.img -no-emul-boot -isohybrid-gpt-basdat -o %s %s/
			`, isoPath, event.Os, isoPath, customIsoPath, isoPath, customIsoPath, preseed, customIsoPath, isoPath, xorrisoIsoPath, customIsoPath)

			command2 := fmt.Sprintf("qemu-img create -f qcow2 %s 15G", diskPath)

			exec.Command(command1)
			exec.Command(command2)

			SendWs(conn, "generating vm xml file")

			xml := CreateVmXML(event.Hostname, diskPath, xorrisoIsoPath, Currentip, Currentport)

			if err := os.WriteFile(xmlPath, []byte(xml), 0644); err != nil {
				fmt.Println("couldnt write xml file", err)
			}

			SendWs(conn, "starting vm")

			command3 := fmt.Sprintf("sudo virsh define %s; sudo virsh start %s", xmlPath, event.Hostname)
			exec.Command(command3)

			SendWs(conn, "done")
		}
	}
}

func CreateVmXML(hostname string, diskpath string, isopath string, ip string, port string) string {
	xml := fmt.Sprintf(`<domain type='kvm'>
	<name>%s</name>
	<memory unit='KiB'>524288</memory>
	<vcpu placement='static'>1</vcpu>
	<os>
		<type arch='x86_64' machine='pc-i440fx-2.12'>hvm</type>
		<boot dev='hd'/>
	</os>
	<devices>
		<!-- Disk -->
		<disk type='file' device='disk'>
		<driver name='qemu' type='qcow2'/>
		<source file='%s'/>
		<target dev='vda' bus='virtio'/>
		<address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x0'/>
		</disk>

		<!-- Network interface -->
		<interface type='network'>
		<source network='default'/>
		<model type='virtio'/>
		<address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
		</interface>

		<disk type='file' device='cdrom'>
			<driver name='qemu' type='raw'/>
			<source file='%s'/>
			<target dev='sdb' bus='sata'/>
			<readonly/>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x05' function='0x0'/>
		</disk>


		<!-- Network interface with a static IP -->
		<interface type='network'>
		<source network='default'/>
		<model type='virtio'/>
		<address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		<protocol family='ipv4'>
			<ip address='%s' prefix='24'/>
		</protocol>
		</interface>
	</devices>

	<!-- Define the qemu:commandline to set the SSH port -->
	<qemu:commandline>
		<qemu:arg value='-netdev'/>
		<qemu:arg value='user,id=usernet,hostfwd=tcp::%s-:22'/>
	</qemu:commandline>
	</domain>
	`, hostname, diskpath, isopath, ip, port)

	return xml
}

func CreatePreseed(password string, hostname string) string {
	preseed := fmt.Sprintf(`# preseed.cfg

	# Choose the language for the installation process
	d-i debian-installer/locale string en_US
	
	# Select the language to be used once installed
	d-i debian-installer/language string en
	
	# Keyboard layout
	d-i keyboard-configuration/xkb-keymap select us
	
	# Set the system clock to UTC
	d-i clock-setup/utc boolean true
	
	# Specify your time zone
	d-i time/zone string Canada/Eastern
	
	# Partitioning
	d-i partman-auto/method string regular
	d-i partman-auto/purge_lvm_from_device boolean true
	d-i partman-lvm/device_remove_lvm boolean true
	d-i partman-md/device_remove_md boolean true
	d-i partman-auto/choose_recipe select atomic
	d-i partman/default_filesystem string ext4
	d-i partman-auto/expert_recipe string \
	   boot-root :: \
	   500 500 500 ext4 \
		  $primary{ } \
		  $bootable{ } \
		  method{ format } \
		  format{ } \
		  use_filesystem{ } \
		  filesystem{ ext4 } \
		  mountpoint{ /boot } \
	   . \
	   14500 14500 14500 ext4 \
		  method{ format } \
		  format{ } \
		  use_filesystem{ } \
		  filesystem{ ext4 } \
		  mountpoint{ / } \
	   . \
	   1000 1000 1000 linux-swap \
		  method{ swap } \
		  format{ } \
	.
	
	# Set the root password
	d-i passwd/root-password password %s
	d-i passwd/root-password-again password %s
	
	# Create a regular user
	d-i passwd/user-fullname string %s
	d-i passwd/username string %s
	d-i passwd/user-password password %s
	d-i passwd/user-password-again password %s
	
	# Setup the bootloader
	d-i grub-installer/only_debian boolean true
	d-i grub-installer/with_other_os boolean true
	d-i finish-install/reboot_in_progress note`, password, password, hostname, hostname, password, password)

	return preseed
}
