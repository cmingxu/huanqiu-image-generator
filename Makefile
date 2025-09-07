
sync:
	rsync -avz . --exclude=node_modules --exclude=.git --exclude=./build dev:/home/xcm/Desktop/huanqiu-image-generator

nextjs:
	ssh dev 'chown -R xcm:xcm /home/xcm/Desktop/huanqiu-image-generator'
	ssh dev 'cd /home/xcm/Desktop/huanqiu-image-generator && docker-compose down && docker-compose up -d --build'
