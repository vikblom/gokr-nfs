.PHONY: mount umount

mount:
	sudo mount -o port=2049,mountport=2049,nfsvers=3,noacl,tcp -t nfs gokrazy.local:/ /run/media/viktor/nfs/

umount:
	sudo umount /run/media/viktor/nfs
