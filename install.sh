n="sillyGirl"
s="/usr/local/$n"
a=arm64
if [[ $(uname -a | grep "x86_64") != "" ]]; then 
    a=amd64
fi ;
if [ ! -d $s ]; then 
    mkdir $s
fi ;
cd $s;
rm -rf $n;
v=`curl https://raw.githubusercontent.com/cdle/binary/main/compile_time.go --silent | tr -cd "[0-9]"`
if [ ${#v} == 13 ]; then
    d="https://raw.githubusercontent.com/cdle/binary/main/sillyGirl_linux_${a}_${v}"
else
    echo "Sorry，你网不好，请使用其他方式下载！"
    exit
fi
echo "检测到版本 $v"
echo "正在从 $d 下载..."
curl -o $n $d && chmod 777 $n
echo "傻妞已安装到 $s"
echo "请手动运行 $s/$n 带 -t 进入交互模式"
