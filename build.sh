# 工作目录
dir="/home/zxh/pomelo_go"

if [ ! -d $dir ]; then
  mkdir $dir
fi

# 进入目录
cd $dir
# 关闭
pkill bench
# 下载
wget -N https://learncloud-beta-enterprise.oss-cn-beijing.aliyuncs.com/library/zhengjiaming/bench/bench.zip
# 解压
unzip -o bench.zip
# 授权
chmod +x bench
# 启动
nohup ./bench &
