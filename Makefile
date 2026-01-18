.PHONY: mount umount

home=$$HOME

mount:
	sudo mount -o port=2049,mountport=2049,nfsvers=3,noacl,tcp -t nfs gokrazy.local:/ $(home)/mnt/nfs/

umount:
	sudo umount $(home)/mnt/nfs
