设计思路：
基于Go编写一个服务，提供一个API，该api无请求参数，返回一个文本数据。主要逻辑如下：
1.访问api时，触发后台的系统爬虫功能；
2.主要时爬取查询数据，通过一个url传递不同参数，且参数是通过cookie方式和session id一起传递的；
3.请求返回一个html页面，页面是一个列表，列表下有查询的总数据，即每页多少项，共多少页，如果没数据，就显示暂时没有记录；
4.通过解析页面，获取上述信息来计算总数据；
5.请求前需要处理登录信息，生成zentaosid，zentaosid是一个32位的uuid，每次触发api请求时随机生成，并使用账号和密码登录；
6.将抓取到的数据统一组织，并以文本方式返回给api请求；
7.api地址：/api/feedback/overview;


AI，待处理：lang=zh-cn; device=desktop; theme=default; feedbackView=0; lastTaskModule=0; lastProduct=2; ajax_lastNext=on; checkedItem=; preProjectID=713; lastProject=734; storyPreProjectID=734; zentaosid=2722e5708efb1a69080faf7a1f7494a9; windowWidth=1912; windowHeight=956
AI，处理中：lang=zh-cn; device=desktop; theme=default; feedbackView=0; lastTaskModule=0; lastProduct=2; ajax_lastNext=on; checkedItem=; preProjectID=713; lastProject=734; storyPreProjectID=734; zentaosid=2722e5708efb1a69080faf7a1f7494a9; windowWidth=1912; windowHeight=956
AI，已处理：lang=zh-cn; device=desktop; theme=default; feedbackView=0; lastTaskModule=0; lastProduct=2; ajax_lastNext=on; checkedItem=; preProjectID=713; lastProject=734; storyPreProjectID=734; zentaosid=2722e5708efb1a69080faf7a1f7494a9; windowWidth=1912; windowHeight=956