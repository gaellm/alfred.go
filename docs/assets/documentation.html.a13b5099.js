import{_ as i,r,o as p,c as u,a as e,b as a,w as n,F as d,d as c,e as s}from"./app.24490de3.js";const h={},k=c('<h1 id="documentation" tabindex="-1"><a class="header-anchor" href="#documentation" aria-hidden="true">#</a> Documentation</h1><p><img src="https://raw.githubusercontent.com/gaellm/alfred.go/main/assets/alfred-study.png" height="180"><span style="vertical-align:middle;position:absolute;margin-top:80px;">This documentation will help you go from a total beginner to a seasoned Alfred.go expert!</span></p><h2 id="what-is-alfred-go" tabindex="-1"><a class="header-anchor" href="#what-is-alfred-go" aria-hidden="true">#</a> What is Alfred.go</h2><p>Alfred.go is an open-source mock that makes performance testing easy and productive for engineering teams. The main goal is to provide a simple way to mock your project partner using json, xml or plain text http queries, without developpement knowledges. Alfred.go has been thinked for cloud projects, it can reach a high request throughput, with a minimum resources footprint. Designed for performance, Alfred.go provides observability: a Prometheus exporter for metrics, and tracing (using OpenTelemetry.io).</p><h2 id="key-features" tabindex="-1"><a class="header-anchor" href="#key-features" aria-hidden="true">#</a> Key features</h2><p>Alfred.go is packed with features, which you can learn all about in the documentation. Key features include:</p><ul><li>mock using simple json files</li><li>use inbound request details to build the mock response using helpers</li><li>generate dates string or random fakers using helpers</li><li>trigger asynchronous actions</li><li>use some javascript to add your features</li><li>live patching your mock without restart Alfred.go</li><li>add a response time offset for a few seconds</li><li>live update the log level</li><li>tracing and metrics</li></ul><h2 id="use-cases" tabindex="-1"><a class="header-anchor" href="#use-cases" aria-hidden="true">#</a> Use cases</h2><p>Users are typically Developers or QA Engineers. They use Alfred.go for creating, testing the performance and reliability of APIs, microservices, and websites without impact on partners. Common use cases are:</p><ul><li><strong>Developpement</strong>: create my app without sending request on external APIs which could be expensive</li><li><strong>Test</strong>: check that my app feature works in all use cases</li><li><strong>Load test</strong>: load test my app without generating lot of requests on my external partners</li><li><strong>Chaos and reliability testing</strong>: if an external call take more time, am I still alive ?</li><li>...</li></ul><h2 id="get-started" tabindex="-1"><a class="header-anchor" href="#get-started" aria-hidden="true">#</a> Get Started</h2><h3 id="install" tabindex="-1"><a class="header-anchor" href="#install" aria-hidden="true">#</a> Install</h3><h4 id="using-the-bundle" tabindex="-1"><a class="header-anchor" href="#using-the-bundle" aria-hidden="true">#</a> Using the Bundle</h4>',13),g=s("The "),f={href:"https://github.com/gaellm/alfred.go/releases",target:"_blank",rel:"noopener noreferrer"},m=s("GitHub Releases page"),_=s(" has a standalone bundle for all platforms. After downloading and extracting the archive for your platform, place the alfred.go or alfred.go.exe binary in your PATH to run Alfred.go in a console from any location."),b=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token function"},"mkdir"),s(" alfred "),e("span",{class:"token operator"},"&&"),s(),e("span",{class:"token function"},"curl"),s(),e("span",{class:"token variable"},[e("span",{class:"token variable"},"$("),e("span",{class:"token function"},"curl"),s(" -s https://api.github.com/repos/gaellm/alfred.go/releases/latest "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(" browser_download_url "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"cut"),s(" -d "),e("span",{class:"token string"},`'"'`),s(" -f "),e("span",{class:"token number"},"4"),s(),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(" `dpkg --print-architecture` "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(" linux"),e("span",{class:"token variable"},")")]),s(" -L "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"tar"),s(" xzv -C alfred "),e("span",{class:"token operator"},"&&"),s(),e("span",{class:"token builtin class-name"},"export"),s(),e("span",{class:"token assign-left variable"},[e("span",{class:"token environment constant"},"PATH")]),e("span",{class:"token operator"},"="),e("span",{class:"token environment constant"},"$PATH"),e("span",{class:"token builtin class-name"},":"),e("span",{class:"token environment constant"},"$PWD"),s(`/alfred/
`)])])],-1),v=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token comment"},"# Prerequisit: brew install dpkg"),s(`
`),e("span",{class:"token function"},"mkdir"),s(" alfred "),e("span",{class:"token operator"},"&&"),s(),e("span",{class:"token function"},"curl"),s(),e("span",{class:"token variable"},[e("span",{class:"token variable"},"$("),e("span",{class:"token function"},"curl"),s(" -s https://api.github.com/repos/gaellm/alfred.go/releases/latest "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(" browser_download_url "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"cut"),s(" -d "),e("span",{class:"token string"},`'"'`),s(" -f "),e("span",{class:"token number"},"4"),s(),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(),e("span",{class:"token punctuation"},"$("),s("dpkg --print-architecture "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"cut"),s(" -d "),e("span",{class:"token string"},'"-"'),s(" -f "),e("span",{class:"token number"},"2"),s(),e("span",{class:"token punctuation"},")"),s(),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"grep"),s(" darwin"),e("span",{class:"token variable"},")")]),s(" -L "),e("span",{class:"token operator"},"|"),s(),e("span",{class:"token function"},"tar"),s(" xzv -C alfred "),e("span",{class:"token operator"},"&&"),s(),e("span",{class:"token builtin class-name"},"export"),s(),e("span",{class:"token assign-left variable"},[e("span",{class:"token environment constant"},"PATH")]),e("span",{class:"token operator"},"="),e("span",{class:"token environment constant"},"$PATH"),e("span",{class:"token builtin class-name"},":"),e("span",{class:"token environment constant"},"$PWD"),s(`/alfred/
`)])])],-1),w=e("div",{class:"language-powershell ext-powershell line-numbers-mode"},[e("pre",{class:"language-powershell"},[e("code",null,[e("span",{class:"token comment"},"# Get the latest release for windows"),s(`
`),e("span",{class:"token variable"},"$bundle"),s(" = "),e("span",{class:"token punctuation"},"("),e("span",{class:"token punctuation"},"("),e("span",{class:"token function"},"Invoke-WebRequest"),s("  "),e("span",{class:"token operator"},"-"),s("UseBasicParsing "),e("span",{class:"token operator"},"-"),s("URI https:"),e("span",{class:"token operator"},"/"),e("span",{class:"token operator"},"/"),s("api"),e("span",{class:"token punctuation"},"."),s("github"),e("span",{class:"token punctuation"},"."),s("com/repos/gaellm/alfred"),e("span",{class:"token punctuation"},"."),s("go/releases/latest"),e("span",{class:"token punctuation"},")"),e("span",{class:"token punctuation"},"."),s("Content "),e("span",{class:"token punctuation"},"|"),s(),e("span",{class:"token function"},"ConvertFrom-Json"),e("span",{class:"token punctuation"},")"),e("span",{class:"token punctuation"},"."),s("assets "),e("span",{class:"token punctuation"},"|"),s(" where "),e("span",{class:"token punctuation"},"{"),s(),e("span",{class:"token variable"},"$_"),e("span",{class:"token punctuation"},"."),s("name "),e("span",{class:"token operator"},"-match"),s(),e("span",{class:"token string"},'"windows"'),s(),e("span",{class:"token punctuation"},"}"),s(),e("span",{class:"token punctuation"},"|"),s(" where "),e("span",{class:"token punctuation"},"{"),s(),e("span",{class:"token variable"},"$_"),e("span",{class:"token punctuation"},"."),s("name "),e("span",{class:"token operator"},"-match"),s(),e("span",{class:"token string"},'"$("'),e("span",{class:"token variable"},"$Env"),s(":PROCESSOR_ARCHITECTURE"),e("span",{class:"token string"},'".ToLower())"'),s(),e("span",{class:"token punctuation"},"}"),s(`
`),e("span",{class:"token comment"},"# download in your curent location and add to path"),s(`
`),e("span",{class:"token function"},"Invoke-WebRequest"),s(),e("span",{class:"token operator"},"-"),s("Uri "),e("span",{class:"token variable"},"$bundle"),e("span",{class:"token punctuation"},"."),s("browser_download_url "),e("span",{class:"token operator"},"-"),s("OutFile "),e("span",{class:"token variable"},"$bundle"),e("span",{class:"token punctuation"},"."),s("name && mkdir alfred && tar "),e("span",{class:"token operator"},"-"),s("xvzf "),e("span",{class:"token variable"},"$bundle"),e("span",{class:"token punctuation"},"."),s("name "),e("span",{class:"token operator"},"-"),s("C alfred && "),e("span",{class:"token variable"},"$env"),s(":Path= "),e("span",{class:"token punctuation"},"("),e("span",{class:"token variable"},"$env"),s(":Path "),e("span",{class:"token operator"},"+"),s(),e("span",{class:"token string"},[s('";'),e("span",{class:"token variable"},"$PWD"),s('/alfred"')]),e("span",{class:"token punctuation"},")"),s(`
`)])]),e("div",{class:"line-numbers","aria-hidden":"true"},[e("span",{class:"line-number"},"1"),e("br"),e("span",{class:"line-number"},"2"),e("br"),e("span",{class:"line-number"},"3"),e("br"),e("span",{class:"line-number"},"4"),e("br")])],-1),y=e("p",null,[s("This archive contains Alfred.go with the default configuration and some examples. If you execute "),e("em",null,"alfred.go"),s(" binary from this folder location, it will load examples mocks files.")],-1),x=e("h4",{id:"docker",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#docker","aria-hidden":"true"},"#"),s(" Docker")],-1),q=s("To use the "),A={href:"https://hub.docker.com/r/gaellm/alfred.go",target:"_blank",rel:"noopener noreferrer"},T=s("Docker image"),C=s(":"),$=c(`<div class="language-console ext-console"><pre class="language-console"><code>$ docker pull gaellm/alfred.go
</code></pre></div><img align="right" src="https://raw.githubusercontent.com/gaellm/alfred.go/main/assets/alfred-docker-style.png" height="180"><p style="margin-top:50px;">This image contains Alfred.go with the default configuration and some examples. It uses a <a href="https://github.com/GoogleContainerTools/distroless">distroless base image</a>, so it contains only the application and its runtime dependencies. The image not contains package managers, shells or any other programs you would expect to find in a standard Linux distribution. Better for security, resources and performance.</p><h3 id="create-a-simple-mock" tabindex="-1"><a class="header-anchor" href="#create-a-simple-mock" aria-hidden="true">#</a> Create a simple mock</h3><p>To do things properly, create a folder named <em>my-alfred-workdir</em>.It will be our sandbox. Inside, create a folder named <em>user-files</em>, this folder will contain a new one named <em>mocks</em>. Here we create a file <em>my-first-mock.json</em> with the following content:</p><div class="language-json ext-json"><pre class="language-json"><code><span class="token punctuation">{</span>
    <span class="token property">&quot;request&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;url&quot;</span><span class="token operator">:</span> <span class="token string">&quot;/my-first-mock&quot;</span>
    <span class="token punctuation">}</span><span class="token punctuation">,</span>
    <span class="token property">&quot;response&quot;</span><span class="token operator">:</span> <span class="token punctuation">{</span>
        <span class="token property">&quot;body&quot;</span><span class="token operator">:</span> <span class="token string">&quot;Congratulations, Sir !&quot;</span>
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre></div><p>Using command lines:</p>`,7),j=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token function"},"mkdir"),s(` -p my-alfred-workdir/user-files/mocks
`),e("span",{class:"token builtin class-name"},"echo"),s(),e("span",{class:"token string"},`'{"request":{"url":"/my-first-mock"},"response":{"body":"Congratulations, Sir !"}}'`),s(),e("span",{class:"token operator"},">"),s(` my-alfred-workdir/user-files/mocks/my-first-mock.json
`)])])],-1),P=e("div",{class:"language-powershell ext-powershell"},[e("pre",{class:"language-powershell"},[e("code",null,[s(`mkdir my-alfred-workdir\\user-files\\mocks
`),e("span",{class:"token function"},"echo"),s(),e("span",{class:"token string"},`'{"request":{"url":"/my-first-mock"},"response":{"body":"Congratulations, Sir !"}}'`),s(" > my-alfred-workdir\\user-files\\mocks\\my-first-mock"),e("span",{class:"token punctuation"},"."),s(`json
`)])])],-1),I=e("p",null,"The directory tree should look like this:",-1),D=e("div",{class:"language-console ext-console"},[e("pre",{class:"language-console"},[e("code",null,`my-alfred-workdir
\u2514\u2500\u2500 user-files
    \u2514\u2500\u2500 mocks
        \u2514\u2500\u2500 my-first-mock.json
`)])],-1),E=e("p",null,[s("Then start Alfred.go from the "),e("em",null,"my-alfred-workdir"),s(" folder:")],-1),G=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,`alfred.go
`)])],-1),H=e("div",{class:"language-powershell ext-powershell"},[e("pre",{class:"language-powershell"},[e("code",null,[s("alfred"),e("span",{class:"token punctuation"},"."),s("go"),e("span",{class:"token punctuation"},"."),s(`exe
`)])])],-1),R=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token function"},"docker"),s(" run --rm --name alfred.go "),e("span",{class:"token punctuation"},"\\"),s(`
  -p `),e("span",{class:"token number"},"8080"),s(":8080 "),e("span",{class:"token punctuation"},"\\"),s(`
  -v `),e("span",{class:"token environment constant"},"$PWD"),s("/user-files/mocks/:/alfred/user-files/mocks/ "),e("span",{class:"token punctuation"},"\\"),s(`
  gaellm/alfred.go
`)])])],-1),U=s("Now access "),W={href:"http://localhost:8080/my-first-mock",target:"_blank",rel:"noopener noreferrer"},B=s("http://localhost:8080/my-first-mock"),L=s(" from your browser. Congratulations, your first step with Alfred.go is done."),S={class:"custom-container tip"},N=e("p",{class:"custom-container-title"},"note",-1),F=s("To display Alfred.go logs more human readables, use "),O={href:"https://stedolan.github.io/jq/",target:"_blank",rel:"noopener noreferrer"},V=s("JQ project"),z=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[s("alfred.go "),e("span",{class:"token operator"},"|"),s(` jq
`)])])],-1),J=e("div",{class:"language-powershell ext-powershell"},[e("pre",{class:"language-powershell"},[e("code",null,[s("alfred"),e("span",{class:"token punctuation"},"."),s("go"),e("span",{class:"token punctuation"},"."),s("exe "),e("span",{class:"token punctuation"},"|"),s(` jq
`)])])],-1),K=e("div",{class:"language-bash ext-sh"},[e("pre",{class:"language-bash"},[e("code",null,[e("span",{class:"token function"},"docker"),s(" run --rm --name alfred.go "),e("span",{class:"token punctuation"},"\\"),s(`
  -p `),e("span",{class:"token number"},"8080"),s(":8080 "),e("span",{class:"token punctuation"},"\\"),s(`
  -v `),e("span",{class:"token environment constant"},"$PWD"),s("/user-files/mocks/:/alfred/user-files/mocks/ "),e("span",{class:"token punctuation"},"\\"),s(`
  gaellm/alfred.go `),e("span",{class:"token operator"},"|"),s(` jq
`)])])],-1),Q=e("h3",{id:"change-default-configuration-port",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#change-default-configuration-port","aria-hidden":"true"},"#"),s(" Change default configuration port")],-1),M=s("By default Alfred.go uses this config file: "),X={href:"https://github.com/gaellm/alfred.go/blob/main/configs/config.json",target:"_blank",rel:"noopener noreferrer"},Y=s("configs/config.json");function Z(ee,se){const o=r("ExternalLinkIcon"),t=r("CodeGroupItem"),l=r("CodeGroup");return p(),u(d,null,[k,e("p",null,[g,e("a",f,[m,a(o)]),_]),a(l,null,{default:n(()=>[a(t,{title:"linux",active:""},{default:n(()=>[b]),_:1}),a(t,{title:"mac"},{default:n(()=>[v]),_:1}),a(t,{title:"powershell"},{default:n(()=>[w]),_:1})]),_:1}),y,x,e("p",null,[q,e("a",A,[T,a(o)]),C]),$,a(l,null,{default:n(()=>[a(t,{title:"bash",active:""},{default:n(()=>[j]),_:1}),a(t,{title:"powershell"},{default:n(()=>[P]),_:1})]),_:1}),I,D,E,a(l,null,{default:n(()=>[a(t,{title:"bash",active:""},{default:n(()=>[G]),_:1}),a(t,{title:"powershell"},{default:n(()=>[H]),_:1}),a(t,{title:"docker"},{default:n(()=>[R]),_:1})]),_:1}),e("p",null,[U,e("a",W,[B,a(o)]),L]),e("div",S,[N,e("p",null,[F,e("a",O,[V,a(o)])]),a(l,null,{default:n(()=>[a(t,{title:"bash",active:""},{default:n(()=>[z]),_:1}),a(t,{title:"powershell"},{default:n(()=>[J]),_:1}),a(t,{title:"docker"},{default:n(()=>[K]),_:1})]),_:1})]),Q,e("p",null,[M,e("a",X,[Y,a(o)])])],64)}var ne=i(h,[["render",Z],["__file","documentation.html.vue"]]);export{ne as default};