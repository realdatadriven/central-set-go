const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["./GlobalForm-CSjeM2Jy.js","./index-E4hdptqK.js","./index-0u4pCFsU.css","./ObjectEditor-CY2UAPRs.js","./NacoEditor-DolCry2b.js","./NacoEditor-ChdkLimv.css"])))=>i.map(i=>d[i]);
var bn=Object.defineProperty;var kn=(r,e,t)=>e in r?bn(r,e,{enumerable:!0,configurable:!0,writable:!0,value:t}):r[e]=t;var W=(r,e,t)=>kn(r,typeof e!="symbol"?e+"":e,t);import{S as _n,i as wn,s as xn,I as pe,e as Re,a as _,b as B,c as J,t as ze,d as g,aa as _e,f as Ge,g as v,h as k,m as X,k as se,p as ue,q as x,r as fe,o as b,as as yn,u as C,v as K,x as dt,y as zt,z as Tn,A as $n,C as zn,D as He,N as En,T as vn,J as Cn,M as Rn,K as Sn,L as In,at as An,l as Ce,R as Ln,X as Se,j as Et,W as On,F as gt,n as Je,w as pt}from"./index-E4hdptqK.js";import{P as qn}from"./Pagination-DMy4LFeK.js";import{N as Gt}from"./NacoEditor-DolCry2b.js";function mt(){return{async:!1,breaks:!1,extensions:null,gfm:!0,hooks:null,pedantic:!1,renderer:null,silent:!1,tokenizer:null,walkTokens:null}}let qe=mt();function Jt(r){qe=r}const Xt=/[&<>"']/,Mn=new RegExp(Xt.source,"g"),Kt=/[<>"']|&(?!(#\d{1,7}|#[Xx][a-fA-F0-9]{1,6}|\w+);)/,Dn=new RegExp(Kt.source,"g"),Nn={"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"},vt=r=>Nn[r];function xe(r,e){if(e){if(Xt.test(r))return r.replace(Mn,vt)}else if(Kt.test(r))return r.replace(Dn,vt);return r}const Bn=/(^|[^\[])\^/g;function U(r,e){let t=typeof r=="string"?r:r.source;e=e||"";const n={replace:(l,s)=>{let i=typeof s=="string"?s:s.source;return i=i.replace(Bn,"$1"),t=t.replace(l,i),n},getRegex:()=>new RegExp(t,e)};return n}function Ct(r){try{r=encodeURI(r).replace(/%25/g,"%")}catch{return null}return r}const Ve={exec:()=>null};function Rt(r,e){const t=r.replace(/\|/g,(s,i,a)=>{let o=!1,p=i;for(;--p>=0&&a[p]==="\\";)o=!o;return o?"|":" |"}),n=t.split(/ \|/);let l=0;if(n[0].trim()||n.shift(),n.length>0&&!n[n.length-1].trim()&&n.pop(),e)if(n.length>e)n.splice(e);else for(;n.length<e;)n.push("");for(;l<n.length;l++)n[l]=n[l].trim().replace(/\\\|/g,"|");return n}function Ue(r,e,t){const n=r.length;if(n===0)return"";let l=0;for(;l<n&&(r.charAt(n-l-1)===e&&!t);)l++;return r.slice(0,n-l)}function Pn(r,e){if(r.indexOf(e[1])===-1)return-1;let t=0;for(let n=0;n<r.length;n++)if(r[n]==="\\")n++;else if(r[n]===e[0])t++;else if(r[n]===e[1]&&(t--,t<0))return n;return-1}function St(r,e,t,n){const l=e.href,s=e.title?xe(e.title):null,i=r[1].replace(/\\([\[\]])/g,"$1");if(r[0].charAt(0)!=="!"){n.state.inLink=!0;const a={type:"link",raw:t,href:l,title:s,text:i,tokens:n.inlineTokens(i)};return n.state.inLink=!1,a}return{type:"image",raw:t,href:l,title:s,text:xe(i)}}function jn(r,e){const t=r.match(/^(\s+)(?:```)/);if(t===null)return e;const n=t[1];return e.split(`
`).map(l=>{const s=l.match(/^\s+/);if(s===null)return l;const[i]=s;return i.length>=n.length?l.slice(n.length):l}).join(`
`)}class st{constructor(e){W(this,"options");W(this,"rules");W(this,"lexer");this.options=e||qe}space(e){const t=this.rules.block.newline.exec(e);if(t&&t[0].length>0)return{type:"space",raw:t[0]}}code(e){const t=this.rules.block.code.exec(e);if(t){const n=t[0].replace(/^(?: {1,4}| {0,3}\t)/gm,"");return{type:"code",raw:t[0],codeBlockStyle:"indented",text:this.options.pedantic?n:Ue(n,`
`)}}}fences(e){const t=this.rules.block.fences.exec(e);if(t){const n=t[0],l=jn(n,t[3]||"");return{type:"code",raw:n,lang:t[2]?t[2].trim().replace(this.rules.inline.anyPunctuation,"$1"):t[2],text:l}}}heading(e){const t=this.rules.block.heading.exec(e);if(t){let n=t[2].trim();if(/#$/.test(n)){const l=Ue(n,"#");(this.options.pedantic||!l||/ $/.test(l))&&(n=l.trim())}return{type:"heading",raw:t[0],depth:t[1].length,text:n,tokens:this.lexer.inline(n)}}}hr(e){const t=this.rules.block.hr.exec(e);if(t)return{type:"hr",raw:Ue(t[0],`
`)}}blockquote(e){const t=this.rules.block.blockquote.exec(e);if(t){let n=Ue(t[0],`
`).split(`
`),l="",s="";const i=[];for(;n.length>0;){let a=!1;const o=[];let p;for(p=0;p<n.length;p++)if(/^ {0,3}>/.test(n[p]))o.push(n[p]),a=!0;else if(!a)o.push(n[p]);else break;n=n.slice(p);const f=o.join(`
`),h=f.replace(/\n {0,3}((?:=+|-+) *)(?=\n|$)/g,`
    $1`).replace(/^ {0,3}>[ \t]?/gm,"");l=l?`${l}
${f}`:f,s=s?`${s}
${h}`:h;const c=this.lexer.state.top;if(this.lexer.state.top=!0,this.lexer.blockTokens(h,i,!0),this.lexer.state.top=c,n.length===0)break;const d=i[i.length-1];if((d==null?void 0:d.type)==="code")break;if((d==null?void 0:d.type)==="blockquote"){const w=d,$=w.raw+`
`+n.join(`
`),E=this.blockquote($);i[i.length-1]=E,l=l.substring(0,l.length-w.raw.length)+E.raw,s=s.substring(0,s.length-w.text.length)+E.text;break}else if((d==null?void 0:d.type)==="list"){const w=d,$=w.raw+`
`+n.join(`
`),E=this.list($);i[i.length-1]=E,l=l.substring(0,l.length-d.raw.length)+E.raw,s=s.substring(0,s.length-w.raw.length)+E.raw,n=$.substring(i[i.length-1].raw.length).split(`
`);continue}}return{type:"blockquote",raw:l,tokens:i,text:s}}}list(e){let t=this.rules.block.list.exec(e);if(t){let n=t[1].trim();const l=n.length>1,s={type:"list",raw:"",ordered:l,start:l?+n.slice(0,-1):"",loose:!1,items:[]};n=l?`\\d{1,9}\\${n.slice(-1)}`:`\\${n}`,this.options.pedantic&&(n=l?n:"[*+-]");const i=new RegExp(`^( {0,3}${n})((?:[	 ][^\\n]*)?(?:\\n|$))`);let a=!1;for(;e;){let o=!1,p="",f="";if(!(t=i.exec(e))||this.rules.block.hr.test(e))break;p=t[0],e=e.substring(p.length);let h=t[2].split(`
`,1)[0].replace(/^\t+/,F=>" ".repeat(3*F.length)),c=e.split(`
`,1)[0],d=!h.trim(),w=0;if(this.options.pedantic?(w=2,f=h.trimStart()):d?w=t[1].length+1:(w=t[2].search(/[^ ]/),w=w>4?1:w,f=h.slice(w),w+=t[1].length),d&&/^[ \t]*$/.test(c)&&(p+=c+`
`,e=e.substring(c.length+1),o=!0),!o){const F=new RegExp(`^ {0,${Math.min(3,w-1)}}(?:[*+-]|\\d{1,9}[.)])((?:[ 	][^\\n]*)?(?:\\n|$))`),M=new RegExp(`^ {0,${Math.min(3,w-1)}}((?:- *){3,}|(?:_ *){3,}|(?:\\* *){3,})(?:\\n+|$)`),R=new RegExp(`^ {0,${Math.min(3,w-1)}}(?:\`\`\`|~~~)`),N=new RegExp(`^ {0,${Math.min(3,w-1)}}#`),oe=new RegExp(`^ {0,${Math.min(3,w-1)}}<(?:[a-z].*>|!--)`,"i");for(;e;){const P=e.split(`
`,1)[0];let j;if(c=P,this.options.pedantic?(c=c.replace(/^ {1,4}(?=( {4})*[^ ])/g,"  "),j=c):j=c.replace(/\t/g,"    "),R.test(c)||N.test(c)||oe.test(c)||F.test(c)||M.test(c))break;if(j.search(/[^ ]/)>=w||!c.trim())f+=`
`+j.slice(w);else{if(d||h.replace(/\t/g,"    ").search(/[^ ]/)>=4||R.test(h)||N.test(h)||M.test(h))break;f+=`
`+c}!d&&!c.trim()&&(d=!0),p+=P+`
`,e=e.substring(P.length+1),h=j.slice(w)}}s.loose||(a?s.loose=!0:/\n[ \t]*\n[ \t]*$/.test(p)&&(a=!0));let $=null,E;this.options.gfm&&($=/^\[[ xX]\] /.exec(f),$&&(E=$[0]!=="[ ] ",f=f.replace(/^\[[ xX]\] +/,""))),s.items.push({type:"list_item",raw:p,task:!!$,checked:E,loose:!1,text:f,tokens:[]}),s.raw+=p}s.items[s.items.length-1].raw=s.items[s.items.length-1].raw.trimEnd(),s.items[s.items.length-1].text=s.items[s.items.length-1].text.trimEnd(),s.raw=s.raw.trimEnd();for(let o=0;o<s.items.length;o++)if(this.lexer.state.top=!1,s.items[o].tokens=this.lexer.blockTokens(s.items[o].text,[]),!s.loose){const p=s.items[o].tokens.filter(h=>h.type==="space"),f=p.length>0&&p.some(h=>/\n.*\n/.test(h.raw));s.loose=f}if(s.loose)for(let o=0;o<s.items.length;o++)s.items[o].loose=!0;return s}}html(e){const t=this.rules.block.html.exec(e);if(t)return{type:"html",block:!0,raw:t[0],pre:t[1]==="pre"||t[1]==="script"||t[1]==="style",text:t[0]}}def(e){const t=this.rules.block.def.exec(e);if(t){const n=t[1].toLowerCase().replace(/\s+/g," "),l=t[2]?t[2].replace(/^<(.*)>$/,"$1").replace(this.rules.inline.anyPunctuation,"$1"):"",s=t[3]?t[3].substring(1,t[3].length-1).replace(this.rules.inline.anyPunctuation,"$1"):t[3];return{type:"def",tag:n,raw:t[0],href:l,title:s}}}table(e){const t=this.rules.block.table.exec(e);if(!t||!/[:|]/.test(t[2]))return;const n=Rt(t[1]),l=t[2].replace(/^\||\| *$/g,"").split("|"),s=t[3]&&t[3].trim()?t[3].replace(/\n[ \t]*$/,"").split(`
`):[],i={type:"table",raw:t[0],header:[],align:[],rows:[]};if(n.length===l.length){for(const a of l)/^ *-+: *$/.test(a)?i.align.push("right"):/^ *:-+: *$/.test(a)?i.align.push("center"):/^ *:-+ *$/.test(a)?i.align.push("left"):i.align.push(null);for(let a=0;a<n.length;a++)i.header.push({text:n[a],tokens:this.lexer.inline(n[a]),header:!0,align:i.align[a]});for(const a of s)i.rows.push(Rt(a,i.header.length).map((o,p)=>({text:o,tokens:this.lexer.inline(o),header:!1,align:i.align[p]})));return i}}lheading(e){const t=this.rules.block.lheading.exec(e);if(t)return{type:"heading",raw:t[0],depth:t[2].charAt(0)==="="?1:2,text:t[1],tokens:this.lexer.inline(t[1])}}paragraph(e){const t=this.rules.block.paragraph.exec(e);if(t){const n=t[1].charAt(t[1].length-1)===`
`?t[1].slice(0,-1):t[1];return{type:"paragraph",raw:t[0],text:n,tokens:this.lexer.inline(n)}}}text(e){const t=this.rules.block.text.exec(e);if(t)return{type:"text",raw:t[0],text:t[0],tokens:this.lexer.inline(t[0])}}escape(e){const t=this.rules.inline.escape.exec(e);if(t)return{type:"escape",raw:t[0],text:xe(t[1])}}tag(e){const t=this.rules.inline.tag.exec(e);if(t)return!this.lexer.state.inLink&&/^<a /i.test(t[0])?this.lexer.state.inLink=!0:this.lexer.state.inLink&&/^<\/a>/i.test(t[0])&&(this.lexer.state.inLink=!1),!this.lexer.state.inRawBlock&&/^<(pre|code|kbd|script)(\s|>)/i.test(t[0])?this.lexer.state.inRawBlock=!0:this.lexer.state.inRawBlock&&/^<\/(pre|code|kbd|script)(\s|>)/i.test(t[0])&&(this.lexer.state.inRawBlock=!1),{type:"html",raw:t[0],inLink:this.lexer.state.inLink,inRawBlock:this.lexer.state.inRawBlock,block:!1,text:t[0]}}link(e){const t=this.rules.inline.link.exec(e);if(t){const n=t[2].trim();if(!this.options.pedantic&&/^</.test(n)){if(!/>$/.test(n))return;const i=Ue(n.slice(0,-1),"\\");if((n.length-i.length)%2===0)return}else{const i=Pn(t[2],"()");if(i>-1){const o=(t[0].indexOf("!")===0?5:4)+t[1].length+i;t[2]=t[2].substring(0,i),t[0]=t[0].substring(0,o).trim(),t[3]=""}}let l=t[2],s="";if(this.options.pedantic){const i=/^([^'"]*[^\s])\s+(['"])(.*)\2/.exec(l);i&&(l=i[1],s=i[3])}else s=t[3]?t[3].slice(1,-1):"";return l=l.trim(),/^</.test(l)&&(this.options.pedantic&&!/>$/.test(n)?l=l.slice(1):l=l.slice(1,-1)),St(t,{href:l&&l.replace(this.rules.inline.anyPunctuation,"$1"),title:s&&s.replace(this.rules.inline.anyPunctuation,"$1")},t[0],this.lexer)}}reflink(e,t){let n;if((n=this.rules.inline.reflink.exec(e))||(n=this.rules.inline.nolink.exec(e))){const l=(n[2]||n[1]).replace(/\s+/g," "),s=t[l.toLowerCase()];if(!s){const i=n[0].charAt(0);return{type:"text",raw:i,text:i}}return St(n,s,n[0],this.lexer)}}emStrong(e,t,n=""){let l=this.rules.inline.emStrongLDelim.exec(e);if(!l||l[3]&&n.match(/[\p{L}\p{N}]/u))return;if(!(l[1]||l[2]||"")||!n||this.rules.inline.punctuation.exec(n)){const i=[...l[0]].length-1;let a,o,p=i,f=0;const h=l[0][0]==="*"?this.rules.inline.emStrongRDelimAst:this.rules.inline.emStrongRDelimUnd;for(h.lastIndex=0,t=t.slice(-1*e.length+i);(l=h.exec(t))!=null;){if(a=l[1]||l[2]||l[3]||l[4]||l[5]||l[6],!a)continue;if(o=[...a].length,l[3]||l[4]){p+=o;continue}else if((l[5]||l[6])&&i%3&&!((i+o)%3)){f+=o;continue}if(p-=o,p>0)continue;o=Math.min(o,o+p+f);const c=[...l[0]][0].length,d=e.slice(0,i+l.index+c+o);if(Math.min(i,o)%2){const $=d.slice(1,-1);return{type:"em",raw:d,text:$,tokens:this.lexer.inlineTokens($)}}const w=d.slice(2,-2);return{type:"strong",raw:d,text:w,tokens:this.lexer.inlineTokens(w)}}}}codespan(e){const t=this.rules.inline.code.exec(e);if(t){let n=t[2].replace(/\n/g," ");const l=/[^ ]/.test(n),s=/^ /.test(n)&&/ $/.test(n);return l&&s&&(n=n.substring(1,n.length-1)),n=xe(n,!0),{type:"codespan",raw:t[0],text:n}}}br(e){const t=this.rules.inline.br.exec(e);if(t)return{type:"br",raw:t[0]}}del(e){const t=this.rules.inline.del.exec(e);if(t)return{type:"del",raw:t[0],text:t[2],tokens:this.lexer.inlineTokens(t[2])}}autolink(e){const t=this.rules.inline.autolink.exec(e);if(t){let n,l;return t[2]==="@"?(n=xe(t[1]),l="mailto:"+n):(n=xe(t[1]),l=n),{type:"link",raw:t[0],text:n,href:l,tokens:[{type:"text",raw:n,text:n}]}}}url(e){var n;let t;if(t=this.rules.inline.url.exec(e)){let l,s;if(t[2]==="@")l=xe(t[0]),s="mailto:"+l;else{let i;do i=t[0],t[0]=((n=this.rules.inline._backpedal.exec(t[0]))==null?void 0:n[0])??"";while(i!==t[0]);l=xe(t[0]),t[1]==="www."?s="http://"+t[0]:s=t[0]}return{type:"link",raw:t[0],text:l,href:s,tokens:[{type:"text",raw:l,text:l}]}}}inlineText(e){const t=this.rules.inline.text.exec(e);if(t){let n;return this.lexer.state.inRawBlock?n=t[0]:n=xe(t[0]),{type:"text",raw:t[0],text:n}}}}const Zn=/^(?:[ \t]*(?:\n|$))+/,Qn=/^((?: {4}| {0,3}\t)[^\n]+(?:\n(?:[ \t]*(?:\n|$))*)?)+/,Hn=/^ {0,3}(`{3,}(?=[^`\n]*(?:\n|$))|~{3,})([^\n]*)(?:\n|$)(?:|([\s\S]*?)(?:\n|$))(?: {0,3}\1[~`]* *(?=\n|$)|$)/,Xe=/^ {0,3}((?:-[\t ]*){3,}|(?:_[ \t]*){3,}|(?:\*[ \t]*){3,})(?:\n+|$)/,Un=/^ {0,3}(#{1,6})(?=\s|$)(.*)(?:\n+|$)/,Yt=/(?:[*+-]|\d{1,9}[.)])/,en=U(/^(?!bull |blockCode|fences|blockquote|heading|html)((?:.|\n(?!\s*?\n|bull |blockCode|fences|blockquote|heading|html))+?)\n {0,3}(=+|-+) *(?:\n+|$)/).replace(/bull/g,Yt).replace(/blockCode/g,/(?: {4}| {0,3}\t)/).replace(/fences/g,/ {0,3}(?:`{3,}|~{3,})/).replace(/blockquote/g,/ {0,3}>/).replace(/heading/g,/ {0,3}#{1,6}/).replace(/html/g,/ {0,3}<[^\n>]+>\n/).getRegex(),bt=/^([^\n]+(?:\n(?!hr|heading|lheading|blockquote|fences|list|html|table| +\n)[^\n]+)*)/,Fn=/^[^\n]+/,kt=/(?!\s*\])(?:\\.|[^\[\]\\])+/,Vn=U(/^ {0,3}\[(label)\]: *(?:\n[ \t]*)?([^<\s][^\s]*|<.*?>)(?:(?: +(?:\n[ \t]*)?| *\n[ \t]*)(title))? *(?:\n+|$)/).replace("label",kt).replace("title",/(?:"(?:\\"?|[^"\\])*"|'[^'\n]*(?:\n[^'\n]+)*\n?'|\([^()]*\))/).getRegex(),Wn=U(/^( {0,3}bull)([ \t][^\n]+?)?(?:\n|$)/).replace(/bull/g,Yt).getRegex(),ot="address|article|aside|base|basefont|blockquote|body|caption|center|col|colgroup|dd|details|dialog|dir|div|dl|dt|fieldset|figcaption|figure|footer|form|frame|frameset|h[1-6]|head|header|hr|html|iframe|legend|li|link|main|menu|menuitem|meta|nav|noframes|ol|optgroup|option|p|param|search|section|summary|table|tbody|td|tfoot|th|thead|title|tr|track|ul",_t=/<!--(?:-?>|[\s\S]*?(?:-->|$))/,Gn=U("^ {0,3}(?:<(script|pre|style|textarea)[\\s>][\\s\\S]*?(?:</\\1>[^\\n]*\\n+|$)|comment[^\\n]*(\\n+|$)|<\\?[\\s\\S]*?(?:\\?>\\n*|$)|<![A-Z][\\s\\S]*?(?:>\\n*|$)|<!\\[CDATA\\[[\\s\\S]*?(?:\\]\\]>\\n*|$)|</?(tag)(?: +|\\n|/?>)[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|<(?!script|pre|style|textarea)([a-z][\\w-]*)(?:attribute)*? */?>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|</(?!script|pre|style|textarea)[a-z][\\w-]*\\s*>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$))","i").replace("comment",_t).replace("tag",ot).replace("attribute",/ +[a-zA-Z:_][\w.:-]*(?: *= *"[^"\n]*"| *= *'[^'\n]*'| *= *[^\s"'=<>`]+)?/).getRegex(),tn=U(bt).replace("hr",Xe).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("|lheading","").replace("|table","").replace("blockquote"," {0,3}>").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",ot).getRegex(),Jn=U(/^( {0,3}> ?(paragraph|[^\n]*)(?:\n|$))+/).replace("paragraph",tn).getRegex(),wt={blockquote:Jn,code:Qn,def:Vn,fences:Hn,heading:Un,hr:Xe,html:Gn,lheading:en,list:Wn,newline:Zn,paragraph:tn,table:Ve,text:Fn},It=U("^ *([^\\n ].*)\\n {0,3}((?:\\| *)?:?-+:? *(?:\\| *:?-+:? *)*(?:\\| *)?)(?:\\n((?:(?! *\\n|hr|heading|blockquote|code|fences|list|html).*(?:\\n|$))*)\\n*|$)").replace("hr",Xe).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("blockquote"," {0,3}>").replace("code","(?: {4}| {0,3}	)[^\\n]").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",ot).getRegex(),Xn={...wt,table:It,paragraph:U(bt).replace("hr",Xe).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("|lheading","").replace("table",It).replace("blockquote"," {0,3}>").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",ot).getRegex()},Kn={...wt,html:U(`^ *(?:comment *(?:\\n|\\s*$)|<(tag)[\\s\\S]+?</\\1> *(?:\\n{2,}|\\s*$)|<tag(?:"[^"]*"|'[^']*'|\\s[^'"/>\\s]*)*?/?> *(?:\\n{2,}|\\s*$))`).replace("comment",_t).replace(/tag/g,"(?!(?:a|em|strong|small|s|cite|q|dfn|abbr|data|time|code|var|samp|kbd|sub|sup|i|b|u|mark|ruby|rt|rp|bdi|bdo|span|br|wbr|ins|del|img)\\b)\\w+(?!:|[^\\w\\s@]*@)\\b").getRegex(),def:/^ *\[([^\]]+)\]: *<?([^\s>]+)>?(?: +(["(][^\n]+[")]))? *(?:\n+|$)/,heading:/^(#{1,6})(.*)(?:\n+|$)/,fences:Ve,lheading:/^(.+?)\n {0,3}(=+|-+) *(?:\n+|$)/,paragraph:U(bt).replace("hr",Xe).replace("heading",` *#{1,6} *[^
]`).replace("lheading",en).replace("|table","").replace("blockquote"," {0,3}>").replace("|fences","").replace("|list","").replace("|html","").replace("|tag","").getRegex()},nn=/^\\([!"#$%&'()*+,\-./:;<=>?@\[\]\\^_`{|}~])/,Yn=/^(`+)([^`]|[^`][\s\S]*?[^`])\1(?!`)/,ln=/^( {2,}|\\)\n(?!\s*$)/,el=/^(`+|[^`])(?:(?= {2,}\n)|[\s\S]*?(?:(?=[\\<!\[`*_]|\b_|$)|[^ ](?= {2,}\n)))/,Ke="\\p{P}\\p{S}",tl=U(/^((?![*_])[\spunctuation])/,"u").replace(/punctuation/g,Ke).getRegex(),nl=/\[[^[\]]*?\]\((?:\\.|[^\\\(\)]|\((?:\\.|[^\\\(\)])*\))*\)|`[^`]*?`|<[^<>]*?>/g,ll=U(/^(?:\*+(?:((?!\*)[punct])|[^\s*]))|^_+(?:((?!_)[punct])|([^\s_]))/,"u").replace(/punct/g,Ke).getRegex(),sl=U("^[^_*]*?__[^_*]*?\\*[^_*]*?(?=__)|[^*]+(?=[^*])|(?!\\*)[punct](\\*+)(?=[\\s]|$)|[^punct\\s](\\*+)(?!\\*)(?=[punct\\s]|$)|(?!\\*)[punct\\s](\\*+)(?=[^punct\\s])|[\\s](\\*+)(?!\\*)(?=[punct])|(?!\\*)[punct](\\*+)(?!\\*)(?=[punct])|[^punct\\s](\\*+)(?=[^punct\\s])","gu").replace(/punct/g,Ke).getRegex(),il=U("^[^_*]*?\\*\\*[^_*]*?_[^_*]*?(?=\\*\\*)|[^_]+(?=[^_])|(?!_)[punct](_+)(?=[\\s]|$)|[^punct\\s](_+)(?!_)(?=[punct\\s]|$)|(?!_)[punct\\s](_+)(?=[^punct\\s])|[\\s](_+)(?!_)(?=[punct])|(?!_)[punct](_+)(?!_)(?=[punct])","gu").replace(/punct/g,Ke).getRegex(),rl=U(/\\([punct])/,"gu").replace(/punct/g,Ke).getRegex(),ol=U(/^<(scheme:[^\s\x00-\x1f<>]*|email)>/).replace("scheme",/[a-zA-Z][a-zA-Z0-9+.-]{1,31}/).replace("email",/[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+(@)[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+(?![-_])/).getRegex(),al=U(_t).replace("(?:-->|$)","-->").getRegex(),cl=U("^comment|^</[a-zA-Z][\\w:-]*\\s*>|^<[a-zA-Z][\\w-]*(?:attribute)*?\\s*/?>|^<\\?[\\s\\S]*?\\?>|^<![a-zA-Z]+\\s[\\s\\S]*?>|^<!\\[CDATA\\[[\\s\\S]*?\\]\\]>").replace("comment",al).replace("attribute",/\s+[a-zA-Z:_][\w.:-]*(?:\s*=\s*"[^"]*"|\s*=\s*'[^']*'|\s*=\s*[^\s"'=<>`]+)?/).getRegex(),it=/(?:\[(?:\\.|[^\[\]\\])*\]|\\.|`[^`]*`|[^\[\]\\`])*?/,ul=U(/^!?\[(label)\]\(\s*(href)(?:\s+(title))?\s*\)/).replace("label",it).replace("href",/<(?:\\.|[^\n<>\\])+>|[^\s\x00-\x1f]*/).replace("title",/"(?:\\"?|[^"\\])*"|'(?:\\'?|[^'\\])*'|\((?:\\\)?|[^)\\])*\)/).getRegex(),sn=U(/^!?\[(label)\]\[(ref)\]/).replace("label",it).replace("ref",kt).getRegex(),rn=U(/^!?\[(ref)\](?:\[\])?/).replace("ref",kt).getRegex(),fl=U("reflink|nolink(?!\\()","g").replace("reflink",sn).replace("nolink",rn).getRegex(),xt={_backpedal:Ve,anyPunctuation:rl,autolink:ol,blockSkip:nl,br:ln,code:Yn,del:Ve,emStrongLDelim:ll,emStrongRDelimAst:sl,emStrongRDelimUnd:il,escape:nn,link:ul,nolink:rn,punctuation:tl,reflink:sn,reflinkSearch:fl,tag:cl,text:el,url:Ve},pl={...xt,link:U(/^!?\[(label)\]\((.*?)\)/).replace("label",it).getRegex(),reflink:U(/^!?\[(label)\]\s*\[([^\]]*)\]/).replace("label",it).getRegex()},ht={...xt,escape:U(nn).replace("])","~|])").getRegex(),url:U(/^((?:ftp|https?):\/\/|www\.)(?:[a-zA-Z0-9\-]+\.?)+[^\s<]*|^email/,"i").replace("email",/[A-Za-z0-9._+-]+(@)[a-zA-Z0-9-_]+(?:\.[a-zA-Z0-9-_]*[a-zA-Z0-9])+(?![-_])/).getRegex(),_backpedal:/(?:[^?!.,:;*_'"~()&]+|\([^)]*\)|&(?![a-zA-Z0-9]+;$)|[?!.,:;*_'"~)]+(?!$))+/,del:/^(~~?)(?=[^\s~])((?:\\.|[^\\])*?(?:\\.|[^\s~\\]))\1(?=[^~]|$)/,text:/^([`~]+|[^`~])(?:(?= {2,}\n)|(?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)|[\s\S]*?(?:(?=[\\<!\[`*~_]|\b_|https?:\/\/|ftp:\/\/|www\.|$)|[^ ](?= {2,}\n)|[^a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-](?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)))/},hl={...ht,br:U(ln).replace("{2,}","*").getRegex(),text:U(ht.text).replace("\\b_","\\b_| {2,}\\n").replace(/\{2,\}/g,"*").getRegex()},lt={normal:wt,gfm:Xn,pedantic:Kn},Fe={normal:xt,gfm:ht,breaks:hl,pedantic:pl};class ye{constructor(e){W(this,"tokens");W(this,"options");W(this,"state");W(this,"tokenizer");W(this,"inlineQueue");this.tokens=[],this.tokens.links=Object.create(null),this.options=e||qe,this.options.tokenizer=this.options.tokenizer||new st,this.tokenizer=this.options.tokenizer,this.tokenizer.options=this.options,this.tokenizer.lexer=this,this.inlineQueue=[],this.state={inLink:!1,inRawBlock:!1,top:!0};const t={block:lt.normal,inline:Fe.normal};this.options.pedantic?(t.block=lt.pedantic,t.inline=Fe.pedantic):this.options.gfm&&(t.block=lt.gfm,this.options.breaks?t.inline=Fe.breaks:t.inline=Fe.gfm),this.tokenizer.rules=t}static get rules(){return{block:lt,inline:Fe}}static lex(e,t){return new ye(t).lex(e)}static lexInline(e,t){return new ye(t).inlineTokens(e)}lex(e){e=e.replace(/\r\n|\r/g,`
`),this.blockTokens(e,this.tokens);for(let t=0;t<this.inlineQueue.length;t++){const n=this.inlineQueue[t];this.inlineTokens(n.src,n.tokens)}return this.inlineQueue=[],this.tokens}blockTokens(e,t=[],n=!1){this.options.pedantic&&(e=e.replace(/\t/g,"    ").replace(/^ +$/gm,""));let l,s,i;for(;e;)if(!(this.options.extensions&&this.options.extensions.block&&this.options.extensions.block.some(a=>(l=a.call({lexer:this},e,t))?(e=e.substring(l.raw.length),t.push(l),!0):!1))){if(l=this.tokenizer.space(e)){e=e.substring(l.raw.length),l.raw.length===1&&t.length>0?t[t.length-1].raw+=`
`:t.push(l);continue}if(l=this.tokenizer.code(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&(s.type==="paragraph"||s.type==="text")?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l);continue}if(l=this.tokenizer.fences(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.heading(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.hr(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.blockquote(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.list(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.html(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.def(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&(s.type==="paragraph"||s.type==="text")?(s.raw+=`
`+l.raw,s.text+=`
`+l.raw,this.inlineQueue[this.inlineQueue.length-1].src=s.text):this.tokens.links[l.tag]||(this.tokens.links[l.tag]={href:l.href,title:l.title});continue}if(l=this.tokenizer.table(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.lheading(e)){e=e.substring(l.raw.length),t.push(l);continue}if(i=e,this.options.extensions&&this.options.extensions.startBlock){let a=1/0;const o=e.slice(1);let p;this.options.extensions.startBlock.forEach(f=>{p=f.call({lexer:this},o),typeof p=="number"&&p>=0&&(a=Math.min(a,p))}),a<1/0&&a>=0&&(i=e.substring(0,a+1))}if(this.state.top&&(l=this.tokenizer.paragraph(i))){s=t[t.length-1],n&&(s==null?void 0:s.type)==="paragraph"?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue.pop(),this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l),n=i.length!==e.length,e=e.substring(l.raw.length);continue}if(l=this.tokenizer.text(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&s.type==="text"?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue.pop(),this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l);continue}if(e){const a="Infinite loop on byte: "+e.charCodeAt(0);if(this.options.silent){console.error(a);break}else throw new Error(a)}}return this.state.top=!0,t}inline(e,t=[]){return this.inlineQueue.push({src:e,tokens:t}),t}inlineTokens(e,t=[]){let n,l,s,i=e,a,o,p;if(this.tokens.links){const f=Object.keys(this.tokens.links);if(f.length>0)for(;(a=this.tokenizer.rules.inline.reflinkSearch.exec(i))!=null;)f.includes(a[0].slice(a[0].lastIndexOf("[")+1,-1))&&(i=i.slice(0,a.index)+"["+"a".repeat(a[0].length-2)+"]"+i.slice(this.tokenizer.rules.inline.reflinkSearch.lastIndex))}for(;(a=this.tokenizer.rules.inline.blockSkip.exec(i))!=null;)i=i.slice(0,a.index)+"["+"a".repeat(a[0].length-2)+"]"+i.slice(this.tokenizer.rules.inline.blockSkip.lastIndex);for(;(a=this.tokenizer.rules.inline.anyPunctuation.exec(i))!=null;)i=i.slice(0,a.index)+"++"+i.slice(this.tokenizer.rules.inline.anyPunctuation.lastIndex);for(;e;)if(o||(p=""),o=!1,!(this.options.extensions&&this.options.extensions.inline&&this.options.extensions.inline.some(f=>(n=f.call({lexer:this},e,t))?(e=e.substring(n.raw.length),t.push(n),!0):!1))){if(n=this.tokenizer.escape(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.tag(e)){e=e.substring(n.raw.length),l=t[t.length-1],l&&n.type==="text"&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(n=this.tokenizer.link(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.reflink(e,this.tokens.links)){e=e.substring(n.raw.length),l=t[t.length-1],l&&n.type==="text"&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(n=this.tokenizer.emStrong(e,i,p)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.codespan(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.br(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.del(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.autolink(e)){e=e.substring(n.raw.length),t.push(n);continue}if(!this.state.inLink&&(n=this.tokenizer.url(e))){e=e.substring(n.raw.length),t.push(n);continue}if(s=e,this.options.extensions&&this.options.extensions.startInline){let f=1/0;const h=e.slice(1);let c;this.options.extensions.startInline.forEach(d=>{c=d.call({lexer:this},h),typeof c=="number"&&c>=0&&(f=Math.min(f,c))}),f<1/0&&f>=0&&(s=e.substring(0,f+1))}if(n=this.tokenizer.inlineText(s)){e=e.substring(n.raw.length),n.raw.slice(-1)!=="_"&&(p=n.raw.slice(-1)),o=!0,l=t[t.length-1],l&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(e){const f="Infinite loop on byte: "+e.charCodeAt(0);if(this.options.silent){console.error(f);break}else throw new Error(f)}}return t}}class rt{constructor(e){W(this,"options");W(this,"parser");this.options=e||qe}space(e){return""}code({text:e,lang:t,escaped:n}){var i;const l=(i=(t||"").match(/^\S*/))==null?void 0:i[0],s=e.replace(/\n$/,"")+`
`;return l?'<pre><code class="language-'+xe(l)+'">'+(n?s:xe(s,!0))+`</code></pre>
`:"<pre><code>"+(n?s:xe(s,!0))+`</code></pre>
`}blockquote({tokens:e}){return`<blockquote>
${this.parser.parse(e)}</blockquote>
`}html({text:e}){return e}heading({tokens:e,depth:t}){return`<h${t}>${this.parser.parseInline(e)}</h${t}>
`}hr(e){return`<hr>
`}list(e){const t=e.ordered,n=e.start;let l="";for(let a=0;a<e.items.length;a++){const o=e.items[a];l+=this.listitem(o)}const s=t?"ol":"ul",i=t&&n!==1?' start="'+n+'"':"";return"<"+s+i+`>
`+l+"</"+s+`>
`}listitem(e){let t="";if(e.task){const n=this.checkbox({checked:!!e.checked});e.loose?e.tokens.length>0&&e.tokens[0].type==="paragraph"?(e.tokens[0].text=n+" "+e.tokens[0].text,e.tokens[0].tokens&&e.tokens[0].tokens.length>0&&e.tokens[0].tokens[0].type==="text"&&(e.tokens[0].tokens[0].text=n+" "+e.tokens[0].tokens[0].text)):e.tokens.unshift({type:"text",raw:n+" ",text:n+" "}):t+=n+" "}return t+=this.parser.parse(e.tokens,!!e.loose),`<li>${t}</li>
`}checkbox({checked:e}){return"<input "+(e?'checked="" ':"")+'disabled="" type="checkbox">'}paragraph({tokens:e}){return`<p>${this.parser.parseInline(e)}</p>
`}table(e){let t="",n="";for(let s=0;s<e.header.length;s++)n+=this.tablecell(e.header[s]);t+=this.tablerow({text:n});let l="";for(let s=0;s<e.rows.length;s++){const i=e.rows[s];n="";for(let a=0;a<i.length;a++)n+=this.tablecell(i[a]);l+=this.tablerow({text:n})}return l&&(l=`<tbody>${l}</tbody>`),`<table>
<thead>
`+t+`</thead>
`+l+`</table>
`}tablerow({text:e}){return`<tr>
${e}</tr>
`}tablecell(e){const t=this.parser.parseInline(e.tokens),n=e.header?"th":"td";return(e.align?`<${n} align="${e.align}">`:`<${n}>`)+t+`</${n}>
`}strong({tokens:e}){return`<strong>${this.parser.parseInline(e)}</strong>`}em({tokens:e}){return`<em>${this.parser.parseInline(e)}</em>`}codespan({text:e}){return`<code>${e}</code>`}br(e){return"<br>"}del({tokens:e}){return`<del>${this.parser.parseInline(e)}</del>`}link({href:e,title:t,tokens:n}){const l=this.parser.parseInline(n),s=Ct(e);if(s===null)return l;e=s;let i='<a href="'+e+'"';return t&&(i+=' title="'+t+'"'),i+=">"+l+"</a>",i}image({href:e,title:t,text:n}){const l=Ct(e);if(l===null)return n;e=l;let s=`<img src="${e}" alt="${n}"`;return t&&(s+=` title="${t}"`),s+=">",s}text(e){return"tokens"in e&&e.tokens?this.parser.parseInline(e.tokens):e.text}}class yt{strong({text:e}){return e}em({text:e}){return e}codespan({text:e}){return e}del({text:e}){return e}html({text:e}){return e}text({text:e}){return e}link({text:e}){return""+e}image({text:e}){return""+e}br(){return""}}class Te{constructor(e){W(this,"options");W(this,"renderer");W(this,"textRenderer");this.options=e||qe,this.options.renderer=this.options.renderer||new rt,this.renderer=this.options.renderer,this.renderer.options=this.options,this.renderer.parser=this,this.textRenderer=new yt}static parse(e,t){return new Te(t).parse(e)}static parseInline(e,t){return new Te(t).parseInline(e)}parse(e,t=!0){let n="";for(let l=0;l<e.length;l++){const s=e[l];if(this.options.extensions&&this.options.extensions.renderers&&this.options.extensions.renderers[s.type]){const a=s,o=this.options.extensions.renderers[a.type].call({parser:this},a);if(o!==!1||!["space","hr","heading","code","table","blockquote","list","html","paragraph","text"].includes(a.type)){n+=o||"";continue}}const i=s;switch(i.type){case"space":{n+=this.renderer.space(i);continue}case"hr":{n+=this.renderer.hr(i);continue}case"heading":{n+=this.renderer.heading(i);continue}case"code":{n+=this.renderer.code(i);continue}case"table":{n+=this.renderer.table(i);continue}case"blockquote":{n+=this.renderer.blockquote(i);continue}case"list":{n+=this.renderer.list(i);continue}case"html":{n+=this.renderer.html(i);continue}case"paragraph":{n+=this.renderer.paragraph(i);continue}case"text":{let a=i,o=this.renderer.text(a);for(;l+1<e.length&&e[l+1].type==="text";)a=e[++l],o+=`
`+this.renderer.text(a);t?n+=this.renderer.paragraph({type:"paragraph",raw:o,text:o,tokens:[{type:"text",raw:o,text:o}]}):n+=o;continue}default:{const a='Token with "'+i.type+'" type was not found.';if(this.options.silent)return console.error(a),"";throw new Error(a)}}}return n}parseInline(e,t){t=t||this.renderer;let n="";for(let l=0;l<e.length;l++){const s=e[l];if(this.options.extensions&&this.options.extensions.renderers&&this.options.extensions.renderers[s.type]){const a=this.options.extensions.renderers[s.type].call({parser:this},s);if(a!==!1||!["escape","html","link","image","strong","em","codespan","br","del","text"].includes(s.type)){n+=a||"";continue}}const i=s;switch(i.type){case"escape":{n+=t.text(i);break}case"html":{n+=t.html(i);break}case"link":{n+=t.link(i);break}case"image":{n+=t.image(i);break}case"strong":{n+=t.strong(i);break}case"em":{n+=t.em(i);break}case"codespan":{n+=t.codespan(i);break}case"br":{n+=t.br(i);break}case"del":{n+=t.del(i);break}case"text":{n+=t.text(i);break}default:{const a='Token with "'+i.type+'" type was not found.';if(this.options.silent)return console.error(a),"";throw new Error(a)}}}return n}}class We{constructor(e){W(this,"options");W(this,"block");this.options=e||qe}preprocess(e){return e}postprocess(e){return e}processAllTokens(e){return e}provideLexer(){return this.block?ye.lex:ye.lexInline}provideParser(){return this.block?Te.parse:Te.parseInline}}W(We,"passThroughHooks",new Set(["preprocess","postprocess","processAllTokens"]));class dl{constructor(...e){W(this,"defaults",mt());W(this,"options",this.setOptions);W(this,"parse",this.parseMarkdown(!0));W(this,"parseInline",this.parseMarkdown(!1));W(this,"Parser",Te);W(this,"Renderer",rt);W(this,"TextRenderer",yt);W(this,"Lexer",ye);W(this,"Tokenizer",st);W(this,"Hooks",We);this.use(...e)}walkTokens(e,t){var l,s;let n=[];for(const i of e)switch(n=n.concat(t.call(this,i)),i.type){case"table":{const a=i;for(const o of a.header)n=n.concat(this.walkTokens(o.tokens,t));for(const o of a.rows)for(const p of o)n=n.concat(this.walkTokens(p.tokens,t));break}case"list":{const a=i;n=n.concat(this.walkTokens(a.items,t));break}default:{const a=i;(s=(l=this.defaults.extensions)==null?void 0:l.childTokens)!=null&&s[a.type]?this.defaults.extensions.childTokens[a.type].forEach(o=>{const p=a[o].flat(1/0);n=n.concat(this.walkTokens(p,t))}):a.tokens&&(n=n.concat(this.walkTokens(a.tokens,t)))}}return n}use(...e){const t=this.defaults.extensions||{renderers:{},childTokens:{}};return e.forEach(n=>{const l={...n};if(l.async=this.defaults.async||l.async||!1,n.extensions&&(n.extensions.forEach(s=>{if(!s.name)throw new Error("extension name required");if("renderer"in s){const i=t.renderers[s.name];i?t.renderers[s.name]=function(...a){let o=s.renderer.apply(this,a);return o===!1&&(o=i.apply(this,a)),o}:t.renderers[s.name]=s.renderer}if("tokenizer"in s){if(!s.level||s.level!=="block"&&s.level!=="inline")throw new Error("extension level must be 'block' or 'inline'");const i=t[s.level];i?i.unshift(s.tokenizer):t[s.level]=[s.tokenizer],s.start&&(s.level==="block"?t.startBlock?t.startBlock.push(s.start):t.startBlock=[s.start]:s.level==="inline"&&(t.startInline?t.startInline.push(s.start):t.startInline=[s.start]))}"childTokens"in s&&s.childTokens&&(t.childTokens[s.name]=s.childTokens)}),l.extensions=t),n.renderer){const s=this.defaults.renderer||new rt(this.defaults);for(const i in n.renderer){if(!(i in s))throw new Error(`renderer '${i}' does not exist`);if(["options","parser"].includes(i))continue;const a=i,o=n.renderer[a],p=s[a];s[a]=(...f)=>{let h=o.apply(s,f);return h===!1&&(h=p.apply(s,f)),h||""}}l.renderer=s}if(n.tokenizer){const s=this.defaults.tokenizer||new st(this.defaults);for(const i in n.tokenizer){if(!(i in s))throw new Error(`tokenizer '${i}' does not exist`);if(["options","rules","lexer"].includes(i))continue;const a=i,o=n.tokenizer[a],p=s[a];s[a]=(...f)=>{let h=o.apply(s,f);return h===!1&&(h=p.apply(s,f)),h}}l.tokenizer=s}if(n.hooks){const s=this.defaults.hooks||new We;for(const i in n.hooks){if(!(i in s))throw new Error(`hook '${i}' does not exist`);if(["options","block"].includes(i))continue;const a=i,o=n.hooks[a],p=s[a];We.passThroughHooks.has(i)?s[a]=f=>{if(this.defaults.async)return Promise.resolve(o.call(s,f)).then(c=>p.call(s,c));const h=o.call(s,f);return p.call(s,h)}:s[a]=(...f)=>{let h=o.apply(s,f);return h===!1&&(h=p.apply(s,f)),h}}l.hooks=s}if(n.walkTokens){const s=this.defaults.walkTokens,i=n.walkTokens;l.walkTokens=function(a){let o=[];return o.push(i.call(this,a)),s&&(o=o.concat(s.call(this,a))),o}}this.defaults={...this.defaults,...l}}),this}setOptions(e){return this.defaults={...this.defaults,...e},this}lexer(e,t){return ye.lex(e,t??this.defaults)}parser(e,t){return Te.parse(e,t??this.defaults)}parseMarkdown(e){return(n,l)=>{const s={...l},i={...this.defaults,...s},a=this.onError(!!i.silent,!!i.async);if(this.defaults.async===!0&&s.async===!1)return a(new Error("marked(): The async option was set to true by an extension. Remove async: false from the parse options object to return a Promise."));if(typeof n>"u"||n===null)return a(new Error("marked(): input parameter is undefined or null"));if(typeof n!="string")return a(new Error("marked(): input parameter is of type "+Object.prototype.toString.call(n)+", string expected"));i.hooks&&(i.hooks.options=i,i.hooks.block=e);const o=i.hooks?i.hooks.provideLexer():e?ye.lex:ye.lexInline,p=i.hooks?i.hooks.provideParser():e?Te.parse:Te.parseInline;if(i.async)return Promise.resolve(i.hooks?i.hooks.preprocess(n):n).then(f=>o(f,i)).then(f=>i.hooks?i.hooks.processAllTokens(f):f).then(f=>i.walkTokens?Promise.all(this.walkTokens(f,i.walkTokens)).then(()=>f):f).then(f=>p(f,i)).then(f=>i.hooks?i.hooks.postprocess(f):f).catch(a);try{i.hooks&&(n=i.hooks.preprocess(n));let f=o(n,i);i.hooks&&(f=i.hooks.processAllTokens(f)),i.walkTokens&&this.walkTokens(f,i.walkTokens);let h=p(f,i);return i.hooks&&(h=i.hooks.postprocess(h)),h}catch(f){return a(f)}}}onError(e,t){return n=>{if(n.message+=`
Please report this to https://github.com/markedjs/marked.`,e){const l="<p>An error occurred:</p><pre>"+xe(n.message+"",!0)+"</pre>";return t?Promise.resolve(l):l}if(t)return Promise.reject(n);throw n}}}const Oe=new dl;function Q(r,e){return Oe.parse(r,e)}Q.options=Q.setOptions=function(r){return Oe.setOptions(r),Q.defaults=Oe.defaults,Jt(Q.defaults),Q};Q.getDefaults=mt;Q.defaults=qe;Q.use=function(...r){return Oe.use(...r),Q.defaults=Oe.defaults,Jt(Q.defaults),Q};Q.walkTokens=function(r,e){return Oe.walkTokens(r,e)};Q.parseInline=Oe.parseInline;Q.Parser=Te;Q.parser=Te.parse;Q.Renderer=rt;Q.TextRenderer=yt;Q.Lexer=ye;Q.lexer=ye.lex;Q.Tokenizer=st;Q.Hooks=We;Q.parse=Q;Q.options;Q.setOptions;Q.use;Q.walkTokens;Q.parseInline;Te.parse;ye.lex;function At(r,e,t){const n=r.slice();return n[67]=e[t],n[68]=e,n[69]=t,n}function Lt(r,e,t){const n=r.slice();return n[70]=e[t],n}function Ot(r,e,t){const n=r.slice();return n[73]=e[t],n}function qt(r,e,t){const n=r.slice();return n[76]=e[t],n}function gl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"play",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function ml(r){let e;return{c(){e=_("span"),g(e,"class","loading loading-spinner")},m(t,n){v(t,e,n)},i:Se,o:Se,d(t){t&&C(e)}}}function Mt(r){let e,t,n,l,s,i;const a=[kl,bl],o=[];function p(h,c){return h[67].isEditing?0:1}t=p(r),n=o[t]=a[t](r);function f(){return r[38](r[67])}return{c(){e=_("button"),n.c(),g(e,"class","btn btn-ghost btn-xs"),g(e,"title","Toggle edit mode"),_e(e,"active",r[67].isEditing)},m(h,c){v(h,e,c),o[t].m(e,null),l=!0,s||(i=se(e,"click",f),s=!0)},p(h,c){r=h;let d=t;t=p(r),t!==d&&(ue(),x(o[d],1,1,()=>{o[d]=null}),fe(),n=o[t],n||(n=o[t]=a[t](r),n.c()),b(n,1),n.m(e,null)),(!l||c[0]&32)&&_e(e,"active",r[67].isEditing)},i(h){l||(b(n),l=!0)},o(h){x(n),l=!1},d(h){h&&C(e),o[t].d(),s=!1,i()}}}function bl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"pencil",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function kl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"eye",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function _l(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"eye-slash",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function wl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"eye",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function Dt(r){let e,t,n,l,s,i,a,o,p,f,h,c,d,w,$,E,F,M,R,N,oe;const P=[yl,xl],j=[];function z(L,O){return L[67].isExecuting?0:1}t=z(r),n=j[t]=P[t](r);function V(){return r[40](r[67])}const A=[$l,Tl],te=[];function de(L,O){return L[67].showOutput?0:1}a=de(r),o=te[a]=A[a](r);function ne(){return r[41](r[67])}c=new pe({props:{type:"outline",icon:"chevron-up",path:"icons"}});function ae(){return r[42](r[67])}E=new pe({props:{type:"outline",icon:"chevron-down",path:"icons"}});function le(){return r[43](r[67])}function be(){r[44].call(M,r[68],r[69])}return{c(){e=_("button"),n.c(),s=B(),i=_("button"),o.c(),p=B(),f=_("button"),h=_("span"),J(c.$$.fragment),d=B(),w=_("button"),$=_("span"),J(E.$$.fragment),F=B(),M=_("input"),g(e,"class","btn btn-ghost btn-xs"),e.disabled=l=r[67].isExecuting,g(e,"title","Run this cell"),_e(e,"executing",r[67].isExecuting),g(i,"class","btn btn-ghost btn-xs"),g(i,"title","Toggle output visibility"),_e(i,"active",r[67].showOutput),g(h,"class","w-5 h-5"),g(f,"class","btn btn-ghost btn-xs"),g(f,"title","Run all cells above"),g($,"class","w-5 h-5"),g(w,"class","btn btn-ghost btn-xs"),g(w,"title","Run all cells below"),g(M,"type","text"),g(M,"placeholder","Variable Name"),g(M,"class","input input-xs")},m(L,O){v(L,e,O),j[t].m(e,null),v(L,s,O),v(L,i,O),te[a].m(i,null),v(L,p,O),v(L,f,O),k(f,h),X(c,h,null),v(L,d,O),v(L,w,O),k(w,$),X(E,$,null),v(L,F,O),v(L,M,O),Et(M,r[67].name),R=!0,N||(oe=[se(e,"click",V),se(i,"click",ne),se(f,"click",ae),se(w,"click",le),se(M,"input",be)],N=!0)},p(L,O){r=L;let ie=t;t=z(r),t!==ie&&(ue(),x(j[ie],1,1,()=>{j[ie]=null}),fe(),n=j[t],n||(n=j[t]=P[t](r),n.c()),b(n,1),n.m(e,null)),(!R||O[0]&32&&l!==(l=r[67].isExecuting))&&(e.disabled=l),(!R||O[0]&32)&&_e(e,"executing",r[67].isExecuting);let I=a;a=de(r),a!==I&&(ue(),x(te[I],1,1,()=>{te[I]=null}),fe(),o=te[a],o||(o=te[a]=A[a](r),o.c()),b(o,1),o.m(i,null)),(!R||O[0]&32)&&_e(i,"active",r[67].showOutput),O[0]&32&&M.value!==r[67].name&&Et(M,r[67].name)},i(L){R||(b(n),b(o),b(c.$$.fragment,L),b(E.$$.fragment,L),R=!0)},o(L){x(n),x(o),x(c.$$.fragment,L),x(E.$$.fragment,L),R=!1},d(L){L&&(C(e),C(s),C(i),C(p),C(f),C(d),C(w),C(F),C(M)),j[t].d(),te[a].d(),K(c),K(E),N=!1,dt(oe)}}}function xl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"play",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function yl(r){let e;return{c(){e=_("span"),g(e,"class","loading loading-spinner")},m(t,n){v(t,e,n)},i:Se,o:Se,d(t){t&&C(e)}}}function Tl(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"chart-bar",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function $l(r){let e,t,n;return t=new pe({props:{type:"outline",icon:"chart-bar",path:"icons"}}),{c(){e=_("span"),J(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){v(l,e,s),X(t,e,null),n=!0},i(l){n||(b(t.$$.fragment,l),n=!0)},o(l){x(t.$$.fragment,l),n=!1},d(l){l&&C(e),K(t)}}}function zl(r){let e,t,n=r[67].showCode&&Nt(r);return{c(){n&&n.c(),e=gt()},m(l,s){n&&n.m(l,s),v(l,e,s),t=!0},p(l,s){l[67].showCode?n?(n.p(l,s),s[0]&32&&b(n,1)):(n=Nt(l),n.c(),b(n,1),n.m(e.parentNode,e)):n&&(ue(),x(n,1,1,()=>{n=null}),fe())},i(l){t||(b(n),t=!0)},o(l){x(n),t=!1},d(l){l&&C(e),n&&n.d(l)}}}function El(r){let e,t,n,l;const s=[Cl,vl],i=[];function a(o,p){return o[67].isEditing&&o[67].showCode?0:!o[67].isEditing&&o[67].showCode?1:-1}return~(e=a(r))&&(t=i[e]=s[e](r)),{c(){t&&t.c(),n=gt()},m(o,p){~e&&i[e].m(o,p),v(o,n,p),l=!0},p(o,p){let f=e;e=a(o),e===f?~e&&i[e].p(o,p):(t&&(ue(),x(i[f],1,1,()=>{i[f]=null}),fe()),~e?(t=i[e],t?t.p(o,p):(t=i[e]=s[e](o),t.c()),b(t,1),t.m(n.parentNode,n)):t=null)},i(o){l||(b(t),l=!0)},o(o){x(t),l=!1},d(o){o&&C(n),~e&&i[e].d(o)}}}function Nt(r){let e,t,n,l;function s(...i){return r[48](r[67],...i)}return t=new Gt({props:{code:r[67].code,language:r[67].type,minHeight:30,maxHeight:700,onEditStop:s,options:{fontSize:12,wordWrap:"on",lineNumbers:"on",minimap:{enabled:!1},tabSize:4}}}),{c(){e=_("div"),n=_("div"),J(t.$$.fragment),Ge(n,"display","contents"),Ge(n,"--theme","vs-dark"),g(e,"class","cell-content")},m(i,a){v(i,e,a),k(e,n),X(t,n,null),l=!0},p(i,a){r=i;const o={};a[0]&32&&(o.code=r[67].code),a[0]&32&&(o.language=r[67].type),a[0]&32&&(o.onEditStop=s),t.$set(o)},i(i){l||(b(t.$$.fragment,i),l=!0)},o(i){x(t.$$.fragment,i),l=!1},d(i){i&&C(e),K(t)}}}function vl(r){let e,t=r[16](r[67].code)+"",n,l;function s(){return r[47](r[67])}return{c(){e=_("div"),g(e,"class","w-full p-0 markdown-rendered")},m(i,a){v(i,e,a),e.innerHTML=t,n||(l=se(e,"dblclick",s),n=!0)},p(i,a){r=i,a[0]&32&&t!==(t=r[16](r[67].code)+"")&&(e.innerHTML=t)},i:Se,o:Se,d(i){i&&C(e),n=!1,l()}}}function Cl(r){let e,t,n,l;function s(...i){return r[46](r[67],...i)}return t=new Gt({props:{code:r[67].code,language:r[67].type,minHeight:30,maxHeight:700,onEditStop:s,options:{fontSize:12,wordWrap:"on",lineNumbers:"on",minimap:{enabled:!1},tabSize:4}}}),{c(){e=_("div"),n=_("div"),J(t.$$.fragment),Ge(n,"display","contents"),Ge(n,"--theme","vs-dark"),g(e,"class","cell-content")},m(i,a){v(i,e,a),k(e,n),X(t,n,null),l=!0},p(i,a){r=i;const o={};a[0]&32&&(o.code=r[67].code),a[0]&32&&(o.language=r[67].type),a[0]&32&&(o.onEditStop=s),t.$set(o)},i(i){l||(b(t.$$.fragment,i),l=!0)},o(i){x(t.$$.fragment,i),l=!1},d(i){i&&C(e),K(t)}}}function Bt(r){let e,t,n,l;const s=[Sl,Rl],i=[];function a(o,p){return o[67].output.error?0:o[67].output.executionTime?1:-1}return~(t=a(r))&&(n=i[t]=s[t](r)),{c(){e=_("div"),n&&n.c(),g(e,"class","cell-output")},m(o,p){v(o,e,p),~t&&i[t].m(e,null),l=!0},p(o,p){let f=t;t=a(o),t===f?~t&&i[t].p(o,p):(n&&(ue(),x(i[f],1,1,()=>{i[f]=null}),fe()),~t?(n=i[t],n?n.p(o,p):(n=i[t]=s[t](o),n.c()),b(n,1),n.m(e,null)):n=null)},i(o){l||(b(n),l=!0)},o(o){x(n),l=!1},d(o){o&&C(e),~t&&i[t].d()}}}function Rl(r){var w;let e,t,n,l,s=r[67].output.executionTime+"",i,a,o,p,f,h,c=((w=r[67].output)==null?void 0:w.data)&&Pt(r),d=r[67].output.data&&jt(r);return{c(){e=_("div"),t=_("span"),t.textContent="✓ Query executed successfully",n=B(),l=_("span"),i=ze(s),a=ze("ms"),o=B(),c&&c.c(),p=B(),d&&d.c(),f=gt(),g(l,"class","execution-time"),g(e,"class","alert alert-success text-sm p-1")},m($,E){v($,e,E),k(e,t),k(e,n),k(e,l),k(l,i),k(l,a),k(e,o),c&&c.m(e,null),v($,p,E),d&&d.m($,E),v($,f,E),h=!0},p($,E){var F;(!h||E[0]&32)&&s!==(s=$[67].output.executionTime+"")&&Je(i,s),(F=$[67].output)!=null&&F.data?c?(c.p($,E),E[0]&32&&b(c,1)):(c=Pt($),c.c(),b(c,1),c.m(e,null)):c&&(ue(),x(c,1,1,()=>{c=null}),fe()),$[67].output.data?d?(d.p($,E),E[0]&32&&b(d,1)):(d=jt($),d.c(),b(d,1),d.m(f.parentNode,f)):d&&(ue(),x(d,1,1,()=>{d=null}),fe())},i($){h||(b(c),b(d),h=!0)},o($){x(c),x(d),h=!1},d($){$&&(C(e),C(p),C(f)),c&&c.d(),d&&d.d($)}}}function Sl(r){let e,t,n,l=r[67].output.error+"",s,i,a,o=r[67].output.executionTime+"",p,f;return{c(){e=_("div"),t=_("strong"),t.textContent="Error:",n=B(),s=ze(l),i=B(),a=_("span"),p=ze(o),f=ze("ms"),g(a,"class","execution-time"),g(e,"class","alert alert-error text-sm p-1")},m(h,c){v(h,e,c),k(e,t),k(e,n),k(e,s),k(e,i),k(e,a),k(a,p),k(a,f)},p(h,c){c[0]&32&&l!==(l=h[67].output.error+"")&&Je(s,l),c[0]&32&&o!==(o=h[67].output.executionTime+"")&&Je(p,o)},i:Se,o:Se,d(h){h&&C(e)}}}function Pt(r){let e,t,n,l,s,i;n=new pe({props:{type:"outline",icon:"arrow-down-tray",path:"icons"}});function a(){return r[49](r[67])}return{c(){e=_("button"),t=_("span"),J(n.$$.fragment),g(t,"class","w-5 h-5"),g(e,"class","btn btn-ghost btn-xs btn-export"),g(e,"title","Export as CSV")},m(o,p){v(o,e,p),k(e,t),X(n,t,null),l=!0,s||(i=se(e,"click",a),s=!0)},p(o,p){r=o},i(o){l||(b(n.$$.fragment,o),l=!0)},o(o){x(n.$$.fragment,o),l=!1},d(o){o&&C(e),K(n),s=!1,i()}}}function jt(r){var $,E,F,M,R,N,oe,P,j;let e,t,n,l,s,i,a,o,p,f=Re(Object.keys(r[67].output.data[0])),h=[];for(let z=0;z<f.length;z+=1)h[z]=Zt(qt(r,f,z));let c=Re(r[67].output.data.slice(((($=r[67].output)==null?void 0:$.offset)||0)*(((E=r[67].output)==null?void 0:E.limit)||10),(((F=r[67].output)==null?void 0:F.offset)||0)*(((M=r[67].output)==null?void 0:M.limit)||10)+(((R=r[67].output)==null?void 0:R.limit)||10))),d=[];for(let z=0;z<c.length;z+=1)d[z]=Ht(Lt(r,c,z));function w(...z){return r[50](r[67],r[69],...z)}return o=new qn({props:{lang:(N=r[9])==null?void 0:N.lang,conf:(oe=r[9])==null?void 0:oe.conf,rows_per_page:((P=r[67].output)==null?void 0:P.limit)||10,page:((j=r[67].output)==null?void 0:j.offset)||0,total_rows:r[67].output.data.length}}),o.$on("pageChange",w),{c(){e=_("div"),t=_("table"),n=_("thead"),l=_("tr");for(let z=0;z<h.length;z+=1)h[z].c();s=B(),i=_("tbody");for(let z=0;z<d.length;z+=1)d[z].c();a=B(),J(o.$$.fragment),g(t,"class","table table-zebra table-xs"),g(e,"class","overflow-x-auto")},m(z,V){v(z,e,V),k(e,t),k(t,n),k(n,l);for(let A=0;A<h.length;A+=1)h[A]&&h[A].m(l,null);k(t,s),k(t,i);for(let A=0;A<d.length;A+=1)d[A]&&d[A].m(i,null);k(e,a),X(o,e,null),p=!0},p(z,V){var te,de,ne,ae,le,be,L,O,ie;if(r=z,V[0]&32){f=Re(Object.keys(r[67].output.data[0]));let I;for(I=0;I<f.length;I+=1){const ge=qt(r,f,I);h[I]?h[I].p(ge,V):(h[I]=Zt(ge),h[I].c(),h[I].m(l,null))}for(;I<h.length;I+=1)h[I].d(1);h.length=f.length}if(V[0]&32){c=Re(r[67].output.data.slice((((te=r[67].output)==null?void 0:te.offset)||0)*(((de=r[67].output)==null?void 0:de.limit)||10),(((ne=r[67].output)==null?void 0:ne.offset)||0)*(((ae=r[67].output)==null?void 0:ae.limit)||10)+(((le=r[67].output)==null?void 0:le.limit)||10)));let I;for(I=0;I<c.length;I+=1){const ge=Lt(r,c,I);d[I]?d[I].p(ge,V):(d[I]=Ht(ge),d[I].c(),d[I].m(i,null))}for(;I<d.length;I+=1)d[I].d(1);d.length=c.length}const A={};V[0]&512&&(A.lang=(be=r[9])==null?void 0:be.lang),V[0]&512&&(A.conf=(L=r[9])==null?void 0:L.conf),V[0]&32&&(A.rows_per_page=((O=r[67].output)==null?void 0:O.limit)||10),V[0]&32&&(A.page=((ie=r[67].output)==null?void 0:ie.offset)||0),V[0]&32&&(A.total_rows=r[67].output.data.length),o.$set(A)},i(z){p||(b(o.$$.fragment,z),p=!0)},o(z){x(o.$$.fragment,z),p=!1},d(z){z&&C(e),pt(h,z),pt(d,z),K(o)}}}function Zt(r){let e,t=r[76]+"",n;return{c(){e=_("th"),n=ze(t)},m(l,s){v(l,e,s),k(e,n)},p(l,s){s[0]&32&&t!==(t=l[76]+"")&&Je(n,t)},d(l){l&&C(e)}}}function Qt(r){let e,t=r[73]+"",n;return{c(){e=_("td"),n=ze(t)},m(l,s){v(l,e,s),k(e,n)},p(l,s){s[0]&32&&t!==(t=l[73]+"")&&Je(n,t)},d(l){l&&C(e)}}}function Ht(r){let e,t,n=Re(Object.values(r[70])),l=[];for(let s=0;s<n.length;s+=1)l[s]=Qt(Ot(r,n,s));return{c(){e=_("tr");for(let s=0;s<l.length;s+=1)l[s].c();t=B()},m(s,i){v(s,e,i);for(let a=0;a<l.length;a+=1)l[a]&&l[a].m(e,null);k(e,t)},p(s,i){if(i[0]&32){n=Re(Object.values(s[70]));let a;for(a=0;a<n.length;a+=1){const o=Ot(s,n,a);l[a]?l[a].p(o,i):(l[a]=Qt(o),l[a].c(),l[a].m(e,t))}for(;a<l.length;a+=1)l[a].d(1);l.length=n.length}},d(s){s&&C(e),pt(l,s)}}}function Ut(r,e){let t,n,l,s,i,a,o,p,f,h,c,d,w,$,E,F,M,R,N,oe,P,j,z,V,A,te,de,ne,ae,le,be,L,O,ie,I;i=new pe({props:{type:"outline",icon:"arrows-pointing-out",path:"icons"}});function ge(...H){return e[37](e[69],...H)}let Y=e[67].type==="markdown"&&Mt(e);const G=[wl,_l],Z=[];function ee(H,he){return H[67].showCode?0:1}f=ee(e),h=Z[f]=G[f](e);function T(){return e[39](e[67])}let y=e[67].type==="sql"&&Dt(e);E=new pe({props:{type:"outline",icon:"trash",path:"icons"}});function me(){return e[45](e[67])}const Ye=[El,zl],Ee=[];function et(H,he){return H[67].type==="markdown"?0:H[67].type==="sql"?1:-1}~(R=et(e))&&(N=Ee[R]=Ye[R](e));let re=e[67].type==="sql"&&e[67].output&&e[67].showOutput&&Bt(e);A=new pe({props:{type:"outline",icon:"plus",path:"icons"}});function at(){return e[51](e[69])}le=new pe({props:{type:"outline",icon:"plus",path:"icons"}});function ct(){return e[52](e[69])}function ut(...H){return e[53](e[69],...H)}function ft(...H){return e[54](e[69],...H)}return{key:r,first:null,c(){t=_("div"),n=_("div"),l=_("button"),s=_("span"),J(i.$$.fragment),a=B(),Y&&Y.c(),o=B(),p=_("button"),h.c(),c=B(),y&&y.c(),d=B(),w=_("button"),$=_("span"),J(E.$$.fragment),F=B(),M=_("div"),N&&N.c(),oe=B(),re&&re.c(),P=B(),j=_("div"),z=_("button"),V=_("span"),J(A.$$.fragment),te=ze(" Code"),de=B(),ne=_("button"),ae=_("span"),J(le.$$.fragment),be=ze(" Markdown"),L=B(),g(s,"class","w-5 h-5"),g(l,"class","btn btn-ghost btn-xs"),g(l,"draggable","true"),g(l,"role","listitem"),g(l,"tabindex","0"),_e(l,"dragging",e[7]===e[69]),g(p,"class","btn btn-ghost btn-xs"),g(p,"title","Toggle code visibility"),_e(p,"active",e[67].showCode),g($,"class","w-5 h-5"),g(w,"class","btn btn-ghost btn-xs"),g(w,"title","Delete this cell"),g(n,"class","absolute left -top-5 transform --translate-x-1/4 z-10 hidden group-hover:flex group-focus-within:flex border border-base-300 bg-base-100 px-0 py-0 space-x-1"),g(M,"class","space-y-4"),g(V,"class","w-5 h-5"),g(z,"class","btn btn-ghost btn-xs"),g(z,"title","Add SQL Cell"),g(ae,"class","w-5 h-5"),g(ne,"class","btn btn-ghost btn-xs"),g(ne,"title","Add Markdown Cell"),g(j,"class","absolute left-1/2 -bottom-5 transform -translate-x-1/2 z-10 hidden group-hover:flex group-focus-within:flex border border-base-300 bg-base-100 px-0 py-0 space-x-1"),g(t,"class","relative group border border-base-200 rounded-xs p-2 pt-2 mb-4 bg-base-100"),_e(t,"drag-over",e[8]===e[69]),this.first=t},m(H,he){v(H,t,he),k(t,n),k(n,l),k(l,s),X(i,s,null),k(n,a),Y&&Y.m(n,null),k(n,o),k(n,p),Z[f].m(p,null),k(n,c),y&&y.m(n,null),k(n,d),k(n,w),k(w,$),X(E,$,null),k(t,F),k(t,M),~R&&Ee[R].m(M,null),k(t,oe),re&&re.m(t,null),k(t,P),k(t,j),k(j,z),k(z,V),X(A,V,null),k(z,te),k(j,de),k(j,ne),k(ne,ae),X(le,ae,null),k(ne,be),k(t,L),O=!0,ie||(I=[se(l,"dragstart",ge),se(p,"click",T),se(w,"click",me),se(z,"click",at),se(ne,"click",ct),se(t,"drop",ut),se(t,"dragover",ft),se(t,"dragleave",e[26]),se(t,"dragend",e[28])],ie=!0)},p(H,he){e=H,(!O||he[0]&160)&&_e(l,"dragging",e[7]===e[69]),e[67].type==="markdown"?Y?(Y.p(e,he),he[0]&32&&b(Y,1)):(Y=Mt(e),Y.c(),b(Y,1),Y.m(n,o)):Y&&(ue(),x(Y,1,1,()=>{Y=null}),fe());let De=f;f=ee(e),f!==De&&(ue(),x(Z[De],1,1,()=>{Z[De]=null}),fe(),h=Z[f],h||(h=Z[f]=G[f](e),h.c()),b(h,1),h.m(p,null)),(!O||he[0]&32)&&_e(p,"active",e[67].showCode),e[67].type==="sql"?y?(y.p(e,he),he[0]&32&&b(y,1)):(y=Dt(e),y.c(),b(y,1),y.m(n,d)):y&&(ue(),x(y,1,1,()=>{y=null}),fe());let Ne=R;R=et(e),R===Ne?~R&&Ee[R].p(e,he):(N&&(ue(),x(Ee[Ne],1,1,()=>{Ee[Ne]=null}),fe()),~R?(N=Ee[R],N?N.p(e,he):(N=Ee[R]=Ye[R](e),N.c()),b(N,1),N.m(M,null)):N=null),e[67].type==="sql"&&e[67].output&&e[67].showOutput?re?(re.p(e,he),he[0]&32&&b(re,1)):(re=Bt(e),re.c(),b(re,1),re.m(t,P)):re&&(ue(),x(re,1,1,()=>{re=null}),fe()),(!O||he[0]&288)&&_e(t,"drag-over",e[8]===e[69])},i(H){O||(b(i.$$.fragment,H),b(Y),b(h),b(y),b(E.$$.fragment,H),b(N),b(re),b(A.$$.fragment,H),b(le.$$.fragment,H),O=!0)},o(H){x(i.$$.fragment,H),x(Y),x(h),x(y),x(E.$$.fragment,H),x(N),x(re),x(A.$$.fragment,H),x(le.$$.fragment,H),O=!1},d(H){H&&C(t),K(i),Y&&Y.d(),Z[f].d(),y&&y.d(),K(E),~R&&Ee[R].d(),re&&re.d(),K(A),K(le),ie=!1,dt(I)}}}function Ft(r){let e,t;return e=new En({props:{open:r[4]}}),{c(){J(e.$$.fragment)},m(n,l){X(e,n,l),t=!0},p(n,l){const s={};l[0]&16&&(s.open=n[4]),e.$set(s)},i(n){t||(b(e.$$.fragment,n),t=!0)},o(n){x(e.$$.fragment,n),t=!1},d(n){K(e,n)}}}function Vt(r){let e,t;return e=new vn({props:{open:r[2].open,type:r[2].type,msg:r[2].msg,timer:2e4}}),e.$on("dismiss",r[55]),{c(){J(e.$$.fragment)},m(n,l){X(e,n,l),t=!0},p(n,l){const s={};l[0]&4&&(s.open=n[2].open),l[0]&4&&(s.type=n[2].type),l[0]&4&&(s.msg=n[2].msg),e.$set(s)},i(n){t||(b(e.$$.fragment,n),t=!0)},o(n){x(e.$$.fragment,n),t=!1},d(n){K(e,n)}}}function Wt(r){let e,t;const n=[r[3]];let l={};for(let s=0;s<n.length;s+=1)l=Cn(l,n[s]);return e=new Rn({props:l}),e.$on("dismiss",r[56]),{c(){J(e.$$.fragment)},m(s,i){X(e,s,i),t=!0},p(s,i){const a=i[0]&8?Sn(n,[In(s[3])]):{};e.$set(a)},i(s){t||(b(e.$$.fragment,s),t=!0)},o(s){x(e.$$.fragment,s),t=!1},d(s){K(e,s)}}}function Il(r){let e,t,n,l,s,i,a,o,p,f,h,c,d,w,$,E,F,M,R,N,oe,P=[],j=new Map,z,V,A,te,de,ne,ae,le,be,L;const O=[ml,gl],ie=[];function I(T,y){return T[6]?0:1}l=I(r),s=ie[l]=O[l](r),p=new pe({props:{type:"outline",icon:"plus",path:"icons"}}),w=new pe({props:{type:"outline",icon:"plus",path:"icons"}}),R=new pe({props:{type:"outline",icon:"save",path:"icons"}});let ge=Re(r[5].cells);const Y=T=>T[67].id;for(let T=0;T<ge.length;T+=1){let y=At(r,ge,T),me=Y(y);j.set(me,P[T]=Ut(me,y))}let G=r[4]&&Ft(r),Z=r[2].open&&Vt(r),ee=r[3].open&&Wt(r);return{c(){e=_("div"),t=_("div"),n=_("button"),s.c(),i=B(),a=_("button"),o=_("span"),J(p.$$.fragment),f=ze(" Code"),h=B(),c=_("button"),d=_("span"),J(w.$$.fragment),$=ze(" Markdown"),E=B(),F=_("button"),M=_("span"),J(R.$$.fragment),N=B(),oe=_("div");for(let T=0;T<P.length;T+=1)P[T].c();z=B(),V=_("style"),V.textContent=`:global(body) {
      margin: 0;
      padding: 0;
      --font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
        sans-serif;
      --background-color: #0f172a;
      --color: #ffffff;
      min-height: 100vh;
      width: 100vw;
      overflow-x: hidden;
    }

    .notebook-container {
      width: 100%;
      padding: 0.5rem;
      min-height: 100vh;
      box-sizing: border-box;
    }

    .notebook-header {
      margin-bottom: 1rem;
      text-align: center;
      padding: 0.5rem 0;
    }

    .notebook-header h1 {
      font-size: 1.75rem;
      font-weight: bold;
      margin-bottom: 0.25rem;
      background: linear-gradient(135deg, #3b82f6, #10b981);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .notebook-header p {
      font-size: 0.875rem;
      --color: #94a3b8;
      margin: 0;
    }

    .notebook-toolbar {
      display: flex;
      gap: 0.25rem;
      margin-bottom: 1rem;
      flex-wrap: wrap;
      padding: 0 0.25rem;
    }

    .notebook-content {
      display: flex;
      flex-direction: column;
      gap: 0.3rem;
      padding: 0 0.25rem;
    }

    .notebook-cell {
      --background-color: #1e293b;
      border-radius: 0.3rem;
      --border: 1px solid #334155;
      overflow: hidden;
      transition: all 0.3s ease;
      width: 100%;
      max-width: 100%;
    }

    .notebook-cell:hover {
      border-color: #ccc;
    }

    .notebook-cell.dragging {
      opacity: 0.5;
      transform: rotate(2deg);
    }

    .notebook-cell.drag-over {
      border-color: #3b82f6;
      box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
    }

    .cell-toolbar {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 0.5rem;
      --background-color: #334155;
      --border-bottom: 1px solid #ccc;
      --min-height: 2rem;
      width: 100%;
    }

    .cell-info {
      display: flex;
      align-items: center;
      gap: 0.5rem;
    }

    .drag-handle {
      cursor: grab;
      --color: #94a3b8;
      font-size: 0.875rem;
      user-select: none;
    }

    .drag-handle:active {
      cursor: grabbing;
    }

    .cell-number {
      font-weight: bold;
      --color: #e2e8f0;
      font-size: 0.75rem;
    }

    .cell-type {
      padding: 0.125rem 0.375rem;
      border-radius: 0.25rem;
      font-size: 0.625rem;
      font-weight: bold;
      --color: white;
    }

    .execution-time {
      font-size: 0.625rem;
      --color: #94a3b8;
    }

    .cell-actions {
      display: flex;
      gap: 0.125rem;
      flex-wrap: wrap;
    }

    /* Ghost Button Styling */
    .btn btn-ghost btn-xs {
      background: transparent;
      border: 1px solid transparent;
      --color: #94a3b8;
      cursor: pointer;
      padding: 0.375rem;
      border-radius: 0.25rem;
      font-size: 0.875rem;
      transition: all 0.2s ease;
      min-width: 2rem;
      height: 2rem;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .btn btn-ghost btn-xs:hover {
      --background-color: rgba(148, 163, 184, 0.1);
      --color: #e2e8f0;
      border-color: rgba(148, 163, 184, 0.2);
    }

    .btn btn-ghost btn-xs:active {
      --background-color: rgba(148, 163, 184, 0.2);
    }

    .btn btn-ghost btn-xs:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn btn-ghost btn-xs.executing {
      --color: #fbbf24;
    }

    .btn btn-ghost btn-xs.active {
      --background-color: rgba(59, 130, 246, 0.1);
      --color: #3b82f6;
      border-color: rgba(59, 130, 246, 0.2);
    }

    .btn-small {
      padding: 0.25rem;
      font-size: 0.75rem;
      min-width: 1.5rem;
      height: 1.5rem;
    }

    .btn-export {
      padding: 0.25rem 0.5rem;
      font-size: 0.75rem;
    }

    .cell-content {
      padding: 0;
    }

    /* Markdown Rendered Content */
    .markdown-rendered {
      --padding: 0.75rem;
      --background-color: #1f2937;
      cursor: pointer;
      --min-height: 50px;
      width: 100%;
      --text-align: left;
    }

    .markdown-rendered:hover {
      --background-color: #374151;
    }

    .markdown-rendered h1 {
      font-size: 1.5rem;
      font-weight: bold;
      margin: 0.5rem 0;
      --color: #f9fafb;
      text-align: left;
    }

    .markdown-rendered h2 {
      font-size: 1.25rem;
      font-weight: bold;
      margin: 0.5rem 0;
      --color: #f3f4f6;
      text-align: left;
    }

    .markdown-rendered h3 {
      font-size: 1.125rem;
      font-weight: bold;
      margin: 0.5rem 0;
      --color: #e5e7eb;
      text-align: left;
    }

    .markdown-rendered p {
      margin: 0.5rem 0;
      line-height: 1.6;
      --color: #d1d5db;
      text-align: left;
    }

    .markdown-rendered ul,
    .markdown-rendered ol {
      margin: 0.5rem 0;
      padding-left: 1.5rem;
      --color: #d1d5db;
      text-align: left;
    }

    .markdown-rendered li {
      margin: 0.25rem 0;
      text-align: left;
    }

    .markdown-rendered strong {
      font-weight: bold;
      --color: #f9fafb;
    }

    .markdown-rendered em {
      font-style: italic;
      --color: #f3f4f6;
    }

    .markdown-rendered code {
      --background-color: #374151;
      padding: 0.125rem 0.25rem;
      border-radius: 0.25rem;
      font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas,
        "Courier New", monospace;
      font-size: 0.875em;
      --color: #fbbf24;
    }

    .markdown-rendered pre {
      --background-color: #111827;
      padding: 0.75rem;
      border-radius: 0.375rem;
      overflow-x: auto;
      margin: 0.5rem 0;
      text-align: left;
    }

    .markdown-rendered pre code {
      background: none;
      padding: 0;
      --color: #e5e7eb;
    }

    /* Monaco Editor Container */
    .monaco-editor-container {
      width: 100%;
      min-height: 60px;
      max-height: 75vh;
      border: 1px solid #374151;
      border-radius: 0.375rem;
      overflow: hidden;
    }

    /* Remove old code editor styles */
    .cell-editor {
      width: 100%;
      padding: 0.5rem;
      border: none;
      --background-color: #1f2937;
      --color: #ffffff;
      resize: none;
      overflow: hidden;
      font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas,
        "Courier New", monospace;
      font-size: 0.75rem;
      line-height: 1.4;
      transition: height 0.1s ease;
    }

    .cell-editor.auto-grow {
      min-height: 50px;
      max-height: 75vh;
    }

    .cell-editor:focus {
      outline: none;
      box-shadow: 0 0 0 1px rgba(59, 130, 246, 0.3);
    }

    .cell-output {
      padding: 0.5rem;
      border-top: 1px solid #ccc;
    }

    .output-success {
      --background-color: #064e3b;
      border: 1px solid #059669;
      border-radius: 0.375rem;
      padding: 0.5rem;
    }

    .output-error {
      --background-color: #7f1d1d;
      border: 1px solid #dc2626;
      border-radius: 0.375rem;
      padding: 0.5rem;
      --color: #fecaca;
    }

    .output-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 0.5rem;
      flex-wrap: wrap;
      gap: 0.25rem;
      font-size: 0.75rem;
    }

    .data-table-container {
      --background-color: #1f2937;
      border-radius: 0.375rem;
      overflow: hidden;
    }

    .data-table {
      width: 100%;
      border-collapse: collapse;
      font-size: 0.75rem;
    }

    .data-table th {
      --background-color: #374151;
      padding: 0.5rem;
      text-align: left;
      font-weight: 600;
      border-bottom: 1px solid #4b5563;
    }

    .data-table td {
      padding: 0.5rem;
      border-bottom: 1px solid #374151;
    }

    .data-table tr:hover {
      --background-color: #374151;
    }

    .table-footer {
      padding: 0.5rem;
      --background-color: #374151;
      --color: #9ca3af;
      font-size: 0.625rem;
      text-align: center;
    }`,A=B(),G&&G.c(),te=B(),Z&&Z.c(),de=B(),ee&&ee.c(),ne=B(),ae=_("div"),g(n,"class","btn btn-ghost btn-xs"),n.disabled=r[6],g(n,"title","Execute All Cells"),_e(n,"executing",r[6]),g(o,"class","w-5 h-5"),g(a,"class","btn btn-ghost btn-xs"),g(a,"title","Add SQL Cell"),g(d,"class","w-5 h-5"),g(c,"class","btn btn-ghost btn-xs"),g(c,"title","Add Markdown Cell"),g(M,"class","w-5 h-5"),g(F,"class","btn btn-ghost btn-xs"),g(F,"title","Save Notebook"),g(t,"class","notebook-toolbar mt-2"),g(oe,"class","notebook-content"),g(e,"class","--flex --flex-col"),g(ae,"class","input-primary input-lg"),Ge(ae,"display","none")},m(T,y){v(T,e,y),k(e,t),k(t,n),ie[l].m(n,null),k(t,i),k(t,a),k(a,o),X(p,o,null),k(a,f),k(t,h),k(t,c),k(c,d),X(w,d,null),k(c,$),k(t,E),k(t,F),k(F,M),X(R,M,null),k(e,N),k(e,oe);for(let me=0;me<P.length;me+=1)P[me]&&P[me].m(oe,null);k(e,z),k(e,V),v(T,A,y),G&&G.m(T,y),v(T,te,y),Z&&Z.m(T,y),v(T,de,y),ee&&ee.m(T,y),v(T,ne,y),v(T,ae,y),le=!0,be||(L=[se(n,"click",r[21]),se(a,"click",r[34]),se(c,"click",r[35]),se(F,"click",r[36])],be=!0)},p(T,y){let me=l;l=I(T),l!==me&&(ue(),x(ie[me],1,1,()=>{ie[me]=null}),fe(),s=ie[l],s||(s=ie[l]=O[l](T),s.c()),b(s,1),s.m(n,null)),(!le||y[0]&64)&&(n.disabled=T[6]),(!le||y[0]&64)&&_e(n,"executing",T[6]),y[0]&1071119264&&(ge=Re(T[5].cells),ue(),P=yn(P,y,Y,1,T,ge,j,oe,An,Ut,null,At),fe()),T[4]?G?(G.p(T,y),y[0]&16&&b(G,1)):(G=Ft(T),G.c(),b(G,1),G.m(te.parentNode,te)):G&&(ue(),x(G,1,1,()=>{G=null}),fe()),T[2].open?Z?(Z.p(T,y),y[0]&4&&b(Z,1)):(Z=Vt(T),Z.c(),b(Z,1),Z.m(de.parentNode,de)):Z&&(ue(),x(Z,1,1,()=>{Z=null}),fe()),T[3].open?ee?(ee.p(T,y),y[0]&8&&b(ee,1)):(ee=Wt(T),ee.c(),b(ee,1),ee.m(ne.parentNode,ne)):ee&&(ue(),x(ee,1,1,()=>{ee=null}),fe())},i(T){if(!le){b(s),b(p.$$.fragment,T),b(w.$$.fragment,T),b(R.$$.fragment,T);for(let y=0;y<ge.length;y+=1)b(P[y]);b(G),b(Z),b(ee),le=!0}},o(T){x(s),x(p.$$.fragment,T),x(w.$$.fragment,T),x(R.$$.fragment,T);for(let y=0;y<P.length;y+=1)x(P[y]);x(G),x(Z),x(ee),le=!1},d(T){T&&(C(e),C(A),C(te),C(de),C(ne),C(ae)),ie[l].d(),K(p),K(w),K(R);for(let y=0;y<P.length;y+=1)P[y].d();G&&G.d(T),Z&&Z.d(T),ee&&ee.d(T),be=!1,dt(L)}}}function Le(){return`cell_${Date.now()}_${Math.random().toString(36).substr(2,9)}`}function Al(r){var a;const e=(a=r==null?void 0:r.output)==null?void 0:a.data,t=Object.keys(e[0]),n=[t.join(","),...e.map(o=>t.map(p=>`"${o[p]}"`).join(","))].join(`
`),l=new Blob([n],{type:"text/csv"}),s=URL.createObjectURL(l),i=document.createElement("a");i.href=s,i.download=`${(r==null?void 0:r.name)||"notebook_data"}.csv`,i.click(),URL.revokeObjectURL(s)}function Ll(r,e,t){let n,l;zt(r,Tn,u=>t(9,n=u)),zt(r,$n,u=>t(57,l=u));let{table:s}=e,{data:i={notebook:"Untitled",active:!0}}=e,{option:a={}}=e,{action:o={}}=e,p={open:!1,msg:null,type:null},f={open:!1,width:12,heigth:85,header:!0,title:"Dialog Title",header_actions:[],component:null},h=!1,c={name:"Sample Data Analysis Notebook",description:"A sample notebook demonstrating SQL queries and markdown documentation",useDuckDB:!0,cells:[{id:Le(),isExecuted:!1,type:"markdown",code:`# Data Analysis Notebook

This notebook demonstrates how to query and analyze data using SQL.

## Getting Started

You can write **markdown** with *formatting*, lists, and more:

- Execute cells individually
- Run all cells at once
- Export results as CSV
- Save your notebook for later
- Execute them outomatic on a cron job defined on the Admin/Cron Jobs to nb/name_of_the_notebook

Powered by DuckDB`,showCode:!0,showOutput:!0,isEditing:!1},{id:Le(),isExecuted:!1,type:"markdown",code:`# CONNECTING TO A DATA SOURCE
        Here for example we are connecting to the data lake for the exemple on etlx`,showCode:!0,showOutput:!0,isEditing:!1},{id:Le(),isExecuted:!1,type:"sql",code:`INSTALL sqlite;
LOAD sqlite;
ATTACH IF NOT EXISTS 'ducklake:sqlite:database/dl_metadata.sqlite' AS dl (DATA_PATH 'database/dl');`,showCode:!0,showOutput:!0,isEditing:!1},{id:Le(),isExecuted:!1,type:"sql",code:"SELECT * FROM dl.lineitem LIMIT 10;",showCode:!0,showOutput:!0,isEditing:!1},{id:Le(),isExecuted:!1,type:"sql",code:"DETACH dl;",showCode:!0,showOutput:!0,isEditing:!1}]},d=!1,w=null,$=null;function E(u){const m={id:Le(),isExecuted:!1,type:u,code:"",showCode:!0,showOutput:!0,isEditing:u==="sql"};t(5,c.cells=[...c.cells,m],c)}function F(u,m){const S={id:Le(),isExecuted:!1,type:u,code:"",showCode:!0,showOutput:!0,isEditing:u==="sql"};t(5,c.cells=[...c.cells.slice(0,m+1),S,...c.cells.slice(m+1)],c)}function M(u){t(5,c.cells=c.cells.filter(m=>m.id!==u),c)}function R(u,m){const S=c.cells.findIndex(D=>D.id===u);S!==-1&&t(5,c.cells[S].code=m,c)}function N(u){const m=c.cells.findIndex(S=>S.id===u);m!==-1&&t(5,c.cells[m].isEditing=!c.cells[m].isEditing,c)}function oe(u){const m=c.cells.findIndex(S=>S.id===u);m!==-1&&t(5,c.cells[m].isEditing=!0,c)}function P(u){try{return Q(u)}catch(m){return`<p style="color: #ef4444;">Error rendering markdown: ${m.message}</p>`}}function j(u){const m=c.cells.findIndex(S=>S.id===u);m!==-1&&t(5,c.cells[m].showCode=!c.cells[m].showCode,c)}function z(u){const m=c.cells.findIndex(S=>S.id===u);m!==-1&&t(5,c.cells[m].showOutput=!c.cells[m].showOutput,c)}const V=async u=>{var D,we,ce,$e,Ie,Be,Pe,je,Ze,ke,Qe,tt,Me,nt,Tt,$t;const m=performance.now(),S={class:"nb",method:"cells",conf:n==null?void 0:n.conf,token:n==null?void 0:n.token,app:n==null?void 0:n.selected_app,data:{...u}};t(4,h=!0);try{const q=await On(S);if(q.success===!0)for(let ve=0;ve<c.cells.length;ve++){const Ae=(D=c.cells[ve])==null?void 0:D.id;if((we=q==null?void 0:q.data)!=null&&we[Ae]){const dn=He(($e=(ce=q==null?void 0:q.data)==null?void 0:ce[Ae])==null?void 0:$e.start).clone(),gn=He((Be=(Ie=q==null?void 0:q.data)==null?void 0:Ie[Ae])==null?void 0:Be.end).clone(),mn=He.duration(gn.diff(dn)).clone();t(5,c.cells[ve].isExecuted=!0,c),t(5,c.cells[ve].isExecuting=!1,c),t(5,c.cells[ve].output={data:Ce.cloneDeep((je=(Pe=q==null?void 0:q.data)==null?void 0:Pe[Ae])==null?void 0:je.data),limit:((Ze=c.cells[ve].output)==null?void 0:Ze.limit)||10,offset:((ke=c.cells[ve].output)==null?void 0:ke.offset)||0,error:((tt=(Qe=q==null?void 0:q.data)==null?void 0:Qe[Ae])==null?void 0:tt.success)===!0?null:(nt=(Me=q==null?void 0:q.data)==null?void 0:Me[Ae])==null?void 0:nt.msg,executionTime:($t=(Tt=q==null?void 0:q.data)==null?void 0:Tt[Ae])!=null&&$t.end?mn:Math.round(performance.now()-m)},c)}}else t(2,p.open=!0,p),t(2,p.type="error",p),t(2,p.msg=q.msg||q.toString(),p)}catch(q){console.log(q,q==null?void 0:q.message),t(2,p.open=!0,p),t(2,p.type="error",p),t(2,p.msg=`${q==null?void 0:q.message}`,p)}t(4,h=!1)},A=u=>u.map(m=>(delete m.output,m));async function te(u){const m=c.cells.findIndex(ce=>ce.id===u);if(m===-1)return;const S=c.cells[m];S.isExecuting=!0,t(5,c.cells=[...c.cells],c);const D=Ce.cloneDeep(c.cells).filter((ce,$e)=>$e<=m&&(ce==null?void 0:ce.isExecuted)||$e===m),we=performance.now();try{await new Promise(ce=>setTimeout(ce,Math.random()*1e3+500)),S.type==="sql"?await V({cells:A(D)}):S.output={executionTime:Math.round(performance.now()-we)},S.isExecuted=!0}catch(ce){S.output={error:ce.message,executionTime:Math.round(performance.now()-we)}}finally{S.isExecuting=!1,t(5,c.cells=[...c.cells],c)}}async function de(){t(6,d=!0),await V({...Ce.cloneDeep(c),cells:A(Ce.cloneDeep(c==null?void 0:c.cells))}),t(6,d=!1)}async function ne(u){const m=c.cells.findIndex(D=>D.id===u),S=[];for(let D=0;D<m;D++)S.push(Ce.cloneDeep(c.cells[D]));await V({cells:A(S)})}async function ae(u){const m=c.cells.findIndex(D=>D.id===u),S=[];for(let D=m+1;D<c.cells.length;D++)S.push(Ce.cloneDeep(c.cells[D]));await V({cells:A(S)})}function le(u,m){t(7,w=m),u.dataTransfer.effectAllowed="move"}function be(u,m){u.preventDefault(),u.dataTransfer.dropEffect="move",t(8,$=m)}function L(){t(8,$=null)}function O(u,m){if(u.preventDefault(),w!==null&&w!==m){const S=c.cells[w],D=[...c.cells];D.splice(w,1);const we=w<m?m-1:m;D.splice(we,0,S),t(5,c.cells=D,c)}t(7,w=null),t(8,$=null)}function ie(){t(7,w=null),t(8,$=null)}const I=u=>{var m,S,D;u.data.offset||(u.data.offset=0),u.data.limit||(u.data.limit=10),u.detail.type==="next"?u.data.offset=((m=u.data)==null?void 0:m.offset)+1:u.detail.type==="back"?u.data.offset=((S=u.data)==null?void 0:S.offset)-1:u.detail.type==="first"?u.data.offset=0:u.detail.type==="last"?u.data.offset=Math.ceil(((u==null?void 0:u._length)||u.data.limit)/u.data.limit)-1:u.detail.type==="rows_per_page"&&(u.data.limit=((D=u==null?void 0:u.detail)==null?void 0:D.rows_per_page)||10,u.data.offset=0),t(5,c={...c})},ge=u=>{console.log("after_save_form".toUpperCase(),u)},Y=u=>{var m,S,D,we,ce,$e;if((D=(S=(m=n==null?void 0:n.tables)==null?void 0:m[u])==null?void 0:S.custom_form)!=null&&D.config)try{return JSON.parse(($e=(ce=(we=n==null?void 0:n.tables)==null?void 0:we[u])==null?void 0:ce.custom_form)==null?void 0:$e.config)}catch(Ie){return console.log(Ie.message),{}}},G=async u=>{var m,S,D,we,ce,$e,Ie,Be,Pe,je,Ze;try{const ke=Ce.cloneDeep(u.data||{});u.clone===!0&&(ke._original_pk=ke[(S=(m=n==null?void 0:n.tables)==null?void 0:m[u.table])==null?void 0:S.pk],delete ke[(we=(D=n==null?void 0:n.tables)==null?void 0:D[u.table])==null?void 0:we.pk],ke.clone=!0);let Qe=`${l.t("crud.add")} - ${($e=(ce=n==null?void 0:n.tables)==null?void 0:ce[u.table])==null?void 0:$e.comment}`;ke!=null&&ke[(Be=(Ie=n==null?void 0:n.tables)==null?void 0:Ie[u.table])==null?void 0:Be.pk]&&(Qe=`${l.t("crud.edit")} - ${(je=(Pe=n==null?void 0:n.tables)==null?void 0:Pe[u.table])==null?void 0:je.comment}`),t(4,h=!0);const tt=(await Ln(async()=>{const{default:nt}=await import("./GlobalForm-CSjeM2Jy.js");return{default:nt}},__vite__mapDeps([0,1,2,3,4,5]),import.meta.url)).default;t(4,h=!1),t(3,f.actions=null,f),t(3,f.header=!0,f),t(3,f.component=tt,f),t(3,f.msg=null,f),t(3,f.html_msg=null,f),t(3,f.header_actions=[{type:"icon",icon:"cog-8-tooth",name:"form_customization",label:l.t("crud.form_customization")},{type:"icon",icon:"save",name:"save",label:l.t("crud.save"),action:ge}],f),t(3,f.title=Qe,f),t(3,f.data={table:u==null?void 0:u.table,data:ke,aux_data:Ce.cloneDeep(u==null?void 0:u.aux_data),main_table:(u==null?void 0:u.main_table)||s},f);const Me=Y(u==null?void 0:u.table);t(3,f.width=((Ze=Me==null?void 0:Me.layout)==null?void 0:Ze.size)||10,f),t(3,f.heigth=null,f),t(3,f.open=!0,f),t(4,h=!1)}catch(ke){t(4,h=!1),console.log(ke.message),t(2,p.open=!0,p),t(2,p.type="error",p),t(2,p.msg=ke.message,p)}},Z=u=>{t(3,f.actions=null,f),t(3,f.open=!1,f)};zn(async()=>{l.changeLanguage(n==null?void 0:n.lang);try{He.locale(n==null?void 0:n.lang)}catch{He.locale("en-us")}console.log("NOTEBOOK  Mount:",{table:s,data:i,option:a});try{i&&i!=null&&i.notebook_conf&&t(5,c=JSON.parse(i==null?void 0:i.notebook_conf))}catch(u){console.log(u.message)}});const ee=()=>E("sql"),T=()=>E("markdown"),y=()=>G({table:s,data:{...i||{notebook:"Untitled",active:!0},notebook_conf:JSON.stringify({...c,cells:A(c.cells)})}}),me=(u,m)=>le(m,u),Ye=u=>N(u.id),Ee=u=>j(u.id),et=u=>te(u.id),re=u=>z(u.id),at=u=>ne(u.id),ct=u=>ae(u.id);function ut(u,m){u[m].name=this.value,t(5,c)}const ft=u=>M(u.id),H=(u,m)=>{R(u.id,m)},he=u=>oe(u.id),De=(u,m)=>{R(u.id,m)},Ne=u=>Al(u),on=(u,m,S)=>I({detail:S.detail,data:u.output,_index:m,_length:u.output.data.length}),an=u=>F("sql",u),cn=u=>F("markdown",u),un=(u,m)=>O(m,u),fn=(u,m)=>be(m,u),pn=u=>{t(2,p.open=!1,p)},hn=u=>Z(u.detail);return r.$$set=u=>{"table"in u&&t(1,s=u.table),"data"in u&&t(0,i=u.data),"option"in u&&t(32,a=u.option),"action"in u&&t(33,o=u.action)},r.$$.update=()=>{r.$$.dirty[0]&2&&console.log({table:s}),r.$$.dirty[1]&4&&((o==null?void 0:o.name)=="CANCEL"&&(o==null||o.action()),(o==null?void 0:o.name)=="SAVE"&&console.log(o),(o==null?void 0:o.name)=="REFRESH"&&console.log(o),(o==null?void 0:o.name)=="RUN_ALL"&&console.log(o))},[i,s,p,f,h,c,d,w,$,n,E,F,M,R,N,oe,P,j,z,A,te,de,ne,ae,le,be,L,O,ie,I,G,Z,a,o,ee,T,y,me,Ye,Ee,et,re,at,ct,ut,ft,H,he,De,Ne,on,an,cn,un,fn,pn,hn]}class Nl extends _n{constructor(e){super(),wn(this,e,Ll,Il,xn,{table:1,data:0,option:32,action:33},null,[-1,-1,-1])}}export{Nl as default};
