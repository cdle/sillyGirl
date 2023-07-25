const fs = require('fs');
const path = require('path');
const vm = require('vm');

// 读取目录中的所有脚本文件
const dirPath = './scripts';
const scripts = fs.readdirSync(dirPath)
  .filter(file => path.extname(file) === '.js')
  .map(file => fs.readFileSync(path.join(dirPath, file), 'utf8'));

// 创建一个沙盒对象
const sandbox = {
  console: console // 在沙盒中暴露 console 对象
};
const context = vm.createContext(sandbox);



// 预读取和运行所有脚本
scripts.forEach(script => {
  const scriptObj = new vm.Script(script, {
    filename: path.join(dirPath, 'script.js')
  });
  scriptObj.runInContext(context);
});