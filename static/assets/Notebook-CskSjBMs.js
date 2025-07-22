const __vite__mapDeps=(i,m=__vite__mapDeps,d=(m.f||(m.f=["./GlobalForm-CfFOlcE5.js","./index-BnwLXrZk.js","./index-BocCPX8f.css","./pt-DIdiQWSC.js","./ObjectEditor-B1tAhWzu.js"])))=>i.map(i=>d[i]);
var hn=Object.defineProperty;var gn=(o,e,t)=>e in o?hn(o,e,{enumerable:!0,configurable:!0,writable:!0,value:t}):o[e]=t;var P=(o,e,t)=>gn(o,typeof e!="symbol"?e+"":e,t);import{S as mn,i as bn,s as kn,I as ae,e as $e,a as y,b as D,c as G,t as ge,d as g,ak as he,f as We,g as $,h as k,m as J,k as ee,p as oe,q as v,r as ie,o as T,al as _n,u as I,v as X,x as ut,y as _t,z as wn,A as xn,C as yn,D as Be,N as Tn,T as zn,J as En,M as vn,K as Rn,L as Cn,am as Sn,l as Se,R as $n,X as Ie,j as wt,W as In,F as ft,n as Pe,w as at,P as An}from"./index-BnwLXrZk.js";import"./pt-DIdiQWSC.js";import{N as Qt}from"./NacoEditor-Bd61ruAw.js";function pt(){return{async:!1,breaks:!1,extensions:null,gfm:!0,hooks:null,pedantic:!1,renderer:null,silent:!1,tokenizer:null,walkTokens:null}}let Me=pt();function Ht(o){Me=o}const Ut=/[&<>"']/,Ln=new RegExp(Ut.source,"g"),Ft=/[<>"']|&(?!(#\d{1,7}|#[Xx][a-fA-F0-9]{1,6}|\w+);)/,On=new RegExp(Ft.source,"g"),qn={"&":"&amp;","<":"&lt;",">":"&gt;",'"':"&quot;","'":"&#39;"},xt=o=>qn[o];function ke(o,e){if(e){if(Ut.test(o))return o.replace(Ln,xt)}else if(Ft.test(o))return o.replace(On,xt);return o}const Mn=/(^|[^\[])\^/g;function B(o,e){let t=typeof o=="string"?o:o.source;e=e||"";const n={replace:(l,s)=>{let r=typeof s=="string"?s:s.source;return r=r.replace(Mn,"$1"),t=t.replace(l,r),n},getRegex:()=>new RegExp(t,e)};return n}function yt(o){try{o=encodeURI(o).replace(/%25/g,"%")}catch{return null}return o}const Fe={exec:()=>null};function Tt(o,e){const t=o.replace(/\|/g,(s,r,i)=>{let a=!1,f=r;for(;--f>=0&&i[f]==="\\";)a=!a;return a?"|":" |"}),n=t.split(/ \|/);let l=0;if(n[0].trim()||n.shift(),n.length>0&&!n[n.length-1].trim()&&n.pop(),e)if(n.length>e)n.splice(e);else for(;n.length<e;)n.push("");for(;l<n.length;l++)n[l]=n[l].trim().replace(/\\\|/g,"|");return n}function He(o,e,t){const n=o.length;if(n===0)return"";let l=0;for(;l<n&&(o.charAt(n-l-1)===e&&!t);)l++;return o.slice(0,n-l)}function Dn(o,e){if(o.indexOf(e[1])===-1)return-1;let t=0;for(let n=0;n<o.length;n++)if(o[n]==="\\")n++;else if(o[n]===e[0])t++;else if(o[n]===e[1]&&(t--,t<0))return n;return-1}function zt(o,e,t,n){const l=e.href,s=e.title?ke(e.title):null,r=o[1].replace(/\\([\[\]])/g,"$1");if(o[0].charAt(0)!=="!"){n.state.inLink=!0;const i={type:"link",raw:t,href:l,title:s,text:r,tokens:n.inlineTokens(r)};return n.state.inLink=!1,i}return{type:"image",raw:t,href:l,title:s,text:ke(r)}}function Nn(o,e){const t=o.match(/^(\s+)(?:```)/);if(t===null)return e;const n=t[1];return e.split(`
`).map(l=>{const s=l.match(/^\s+/);if(s===null)return l;const[r]=s;return r.length>=n.length?l.slice(n.length):l}).join(`
`)}class et{constructor(e){P(this,"options");P(this,"rules");P(this,"lexer");this.options=e||Me}space(e){const t=this.rules.block.newline.exec(e);if(t&&t[0].length>0)return{type:"space",raw:t[0]}}code(e){const t=this.rules.block.code.exec(e);if(t){const n=t[0].replace(/^(?: {1,4}| {0,3}\t)/gm,"");return{type:"code",raw:t[0],codeBlockStyle:"indented",text:this.options.pedantic?n:He(n,`
`)}}}fences(e){const t=this.rules.block.fences.exec(e);if(t){const n=t[0],l=Nn(n,t[3]||"");return{type:"code",raw:n,lang:t[2]?t[2].trim().replace(this.rules.inline.anyPunctuation,"$1"):t[2],text:l}}}heading(e){const t=this.rules.block.heading.exec(e);if(t){let n=t[2].trim();if(/#$/.test(n)){const l=He(n,"#");(this.options.pedantic||!l||/ $/.test(l))&&(n=l.trim())}return{type:"heading",raw:t[0],depth:t[1].length,text:n,tokens:this.lexer.inline(n)}}}hr(e){const t=this.rules.block.hr.exec(e);if(t)return{type:"hr",raw:He(t[0],`
`)}}blockquote(e){const t=this.rules.block.blockquote.exec(e);if(t){let n=He(t[0],`
`).split(`
`),l="",s="";const r=[];for(;n.length>0;){let i=!1;const a=[];let f;for(f=0;f<n.length;f++)if(/^ {0,3}>/.test(n[f]))a.push(n[f]),i=!0;else if(!i)a.push(n[f]);else break;n=n.slice(f);const u=a.join(`
`),d=u.replace(/\n {0,3}((?:=+|-+) *)(?=\n|$)/g,`
    $1`).replace(/^ {0,3}>[ \t]?/gm,"");l=l?`${l}
${u}`:u,s=s?`${s}
${d}`:d;const c=this.lexer.state.top;if(this.lexer.state.top=!0,this.lexer.blockTokens(d,r,!0),this.lexer.state.top=c,n.length===0)break;const _=r[r.length-1];if((_==null?void 0:_.type)==="code")break;if((_==null?void 0:_.type)==="blockquote"){const h=_,m=h.raw+`
`+n.join(`
`),w=this.blockquote(m);r[r.length-1]=w,l=l.substring(0,l.length-h.raw.length)+w.raw,s=s.substring(0,s.length-h.text.length)+w.text;break}else if((_==null?void 0:_.type)==="list"){const h=_,m=h.raw+`
`+n.join(`
`),w=this.list(m);r[r.length-1]=w,l=l.substring(0,l.length-_.raw.length)+w.raw,s=s.substring(0,s.length-h.raw.length)+w.raw,n=m.substring(r[r.length-1].raw.length).split(`
`);continue}}return{type:"blockquote",raw:l,tokens:r,text:s}}}list(e){let t=this.rules.block.list.exec(e);if(t){let n=t[1].trim();const l=n.length>1,s={type:"list",raw:"",ordered:l,start:l?+n.slice(0,-1):"",loose:!1,items:[]};n=l?`\\d{1,9}\\${n.slice(-1)}`:`\\${n}`,this.options.pedantic&&(n=l?n:"[*+-]");const r=new RegExp(`^( {0,3}${n})((?:[	 ][^\\n]*)?(?:\\n|$))`);let i=!1;for(;e;){let a=!1,f="",u="";if(!(t=r.exec(e))||this.rules.block.hr.test(e))break;f=t[0],e=e.substring(f.length);let d=t[2].split(`
`,1)[0].replace(/^\t+/,R=>" ".repeat(3*R.length)),c=e.split(`
`,1)[0],_=!d.trim(),h=0;if(this.options.pedantic?(h=2,u=d.trimStart()):_?h=t[1].length+1:(h=t[2].search(/[^ ]/),h=h>4?1:h,u=d.slice(h),h+=t[1].length),_&&/^[ \t]*$/.test(c)&&(f+=c+`
`,e=e.substring(c.length+1),a=!0),!a){const R=new RegExp(`^ {0,${Math.min(3,h-1)}}(?:[*+-]|\\d{1,9}[.)])((?:[ 	][^\\n]*)?(?:\\n|$))`),A=new RegExp(`^ {0,${Math.min(3,h-1)}}((?:- *){3,}|(?:_ *){3,}|(?:\\* *){3,})(?:\\n+|$)`),x=new RegExp(`^ {0,${Math.min(3,h-1)}}(?:\`\`\`|~~~)`),j=new RegExp(`^ {0,${Math.min(3,h-1)}}#`),te=new RegExp(`^ {0,${Math.min(3,h-1)}}<(?:[a-z].*>|!--)`,"i");for(;e;){const q=e.split(`
`,1)[0];let W;if(c=q,this.options.pedantic?(c=c.replace(/^ {1,4}(?=( {4})*[^ ])/g,"  "),W=c):W=c.replace(/\t/g,"    "),x.test(c)||j.test(c)||te.test(c)||R.test(c)||A.test(c))break;if(W.search(/[^ ]/)>=h||!c.trim())u+=`
`+W.slice(h);else{if(_||d.replace(/\t/g,"    ").search(/[^ ]/)>=4||x.test(d)||j.test(d)||A.test(d))break;u+=`
`+c}!_&&!c.trim()&&(_=!0),f+=q+`
`,e=e.substring(q.length+1),d=W.slice(h)}}s.loose||(i?s.loose=!0:/\n[ \t]*\n[ \t]*$/.test(f)&&(i=!0));let m=null,w;this.options.gfm&&(m=/^\[[ xX]\] /.exec(u),m&&(w=m[0]!=="[ ] ",u=u.replace(/^\[[ xX]\] +/,""))),s.items.push({type:"list_item",raw:f,task:!!m,checked:w,loose:!1,text:u,tokens:[]}),s.raw+=f}s.items[s.items.length-1].raw=s.items[s.items.length-1].raw.trimEnd(),s.items[s.items.length-1].text=s.items[s.items.length-1].text.trimEnd(),s.raw=s.raw.trimEnd();for(let a=0;a<s.items.length;a++)if(this.lexer.state.top=!1,s.items[a].tokens=this.lexer.blockTokens(s.items[a].text,[]),!s.loose){const f=s.items[a].tokens.filter(d=>d.type==="space"),u=f.length>0&&f.some(d=>/\n.*\n/.test(d.raw));s.loose=u}if(s.loose)for(let a=0;a<s.items.length;a++)s.items[a].loose=!0;return s}}html(e){const t=this.rules.block.html.exec(e);if(t)return{type:"html",block:!0,raw:t[0],pre:t[1]==="pre"||t[1]==="script"||t[1]==="style",text:t[0]}}def(e){const t=this.rules.block.def.exec(e);if(t){const n=t[1].toLowerCase().replace(/\s+/g," "),l=t[2]?t[2].replace(/^<(.*)>$/,"$1").replace(this.rules.inline.anyPunctuation,"$1"):"",s=t[3]?t[3].substring(1,t[3].length-1).replace(this.rules.inline.anyPunctuation,"$1"):t[3];return{type:"def",tag:n,raw:t[0],href:l,title:s}}}table(e){const t=this.rules.block.table.exec(e);if(!t||!/[:|]/.test(t[2]))return;const n=Tt(t[1]),l=t[2].replace(/^\||\| *$/g,"").split("|"),s=t[3]&&t[3].trim()?t[3].replace(/\n[ \t]*$/,"").split(`
`):[],r={type:"table",raw:t[0],header:[],align:[],rows:[]};if(n.length===l.length){for(const i of l)/^ *-+: *$/.test(i)?r.align.push("right"):/^ *:-+: *$/.test(i)?r.align.push("center"):/^ *:-+ *$/.test(i)?r.align.push("left"):r.align.push(null);for(let i=0;i<n.length;i++)r.header.push({text:n[i],tokens:this.lexer.inline(n[i]),header:!0,align:r.align[i]});for(const i of s)r.rows.push(Tt(i,r.header.length).map((a,f)=>({text:a,tokens:this.lexer.inline(a),header:!1,align:r.align[f]})));return r}}lheading(e){const t=this.rules.block.lheading.exec(e);if(t)return{type:"heading",raw:t[0],depth:t[2].charAt(0)==="="?1:2,text:t[1],tokens:this.lexer.inline(t[1])}}paragraph(e){const t=this.rules.block.paragraph.exec(e);if(t){const n=t[1].charAt(t[1].length-1)===`
`?t[1].slice(0,-1):t[1];return{type:"paragraph",raw:t[0],text:n,tokens:this.lexer.inline(n)}}}text(e){const t=this.rules.block.text.exec(e);if(t)return{type:"text",raw:t[0],text:t[0],tokens:this.lexer.inline(t[0])}}escape(e){const t=this.rules.inline.escape.exec(e);if(t)return{type:"escape",raw:t[0],text:ke(t[1])}}tag(e){const t=this.rules.inline.tag.exec(e);if(t)return!this.lexer.state.inLink&&/^<a /i.test(t[0])?this.lexer.state.inLink=!0:this.lexer.state.inLink&&/^<\/a>/i.test(t[0])&&(this.lexer.state.inLink=!1),!this.lexer.state.inRawBlock&&/^<(pre|code|kbd|script)(\s|>)/i.test(t[0])?this.lexer.state.inRawBlock=!0:this.lexer.state.inRawBlock&&/^<\/(pre|code|kbd|script)(\s|>)/i.test(t[0])&&(this.lexer.state.inRawBlock=!1),{type:"html",raw:t[0],inLink:this.lexer.state.inLink,inRawBlock:this.lexer.state.inRawBlock,block:!1,text:t[0]}}link(e){const t=this.rules.inline.link.exec(e);if(t){const n=t[2].trim();if(!this.options.pedantic&&/^</.test(n)){if(!/>$/.test(n))return;const r=He(n.slice(0,-1),"\\");if((n.length-r.length)%2===0)return}else{const r=Dn(t[2],"()");if(r>-1){const a=(t[0].indexOf("!")===0?5:4)+t[1].length+r;t[2]=t[2].substring(0,r),t[0]=t[0].substring(0,a).trim(),t[3]=""}}let l=t[2],s="";if(this.options.pedantic){const r=/^([^'"]*[^\s])\s+(['"])(.*)\2/.exec(l);r&&(l=r[1],s=r[3])}else s=t[3]?t[3].slice(1,-1):"";return l=l.trim(),/^</.test(l)&&(this.options.pedantic&&!/>$/.test(n)?l=l.slice(1):l=l.slice(1,-1)),zt(t,{href:l&&l.replace(this.rules.inline.anyPunctuation,"$1"),title:s&&s.replace(this.rules.inline.anyPunctuation,"$1")},t[0],this.lexer)}}reflink(e,t){let n;if((n=this.rules.inline.reflink.exec(e))||(n=this.rules.inline.nolink.exec(e))){const l=(n[2]||n[1]).replace(/\s+/g," "),s=t[l.toLowerCase()];if(!s){const r=n[0].charAt(0);return{type:"text",raw:r,text:r}}return zt(n,s,n[0],this.lexer)}}emStrong(e,t,n=""){let l=this.rules.inline.emStrongLDelim.exec(e);if(!l||l[3]&&n.match(/[\p{L}\p{N}]/u))return;if(!(l[1]||l[2]||"")||!n||this.rules.inline.punctuation.exec(n)){const r=[...l[0]].length-1;let i,a,f=r,u=0;const d=l[0][0]==="*"?this.rules.inline.emStrongRDelimAst:this.rules.inline.emStrongRDelimUnd;for(d.lastIndex=0,t=t.slice(-1*e.length+r);(l=d.exec(t))!=null;){if(i=l[1]||l[2]||l[3]||l[4]||l[5]||l[6],!i)continue;if(a=[...i].length,l[3]||l[4]){f+=a;continue}else if((l[5]||l[6])&&r%3&&!((r+a)%3)){u+=a;continue}if(f-=a,f>0)continue;a=Math.min(a,a+f+u);const c=[...l[0]][0].length,_=e.slice(0,r+l.index+c+a);if(Math.min(r,a)%2){const m=_.slice(1,-1);return{type:"em",raw:_,text:m,tokens:this.lexer.inlineTokens(m)}}const h=_.slice(2,-2);return{type:"strong",raw:_,text:h,tokens:this.lexer.inlineTokens(h)}}}}codespan(e){const t=this.rules.inline.code.exec(e);if(t){let n=t[2].replace(/\n/g," ");const l=/[^ ]/.test(n),s=/^ /.test(n)&&/ $/.test(n);return l&&s&&(n=n.substring(1,n.length-1)),n=ke(n,!0),{type:"codespan",raw:t[0],text:n}}}br(e){const t=this.rules.inline.br.exec(e);if(t)return{type:"br",raw:t[0]}}del(e){const t=this.rules.inline.del.exec(e);if(t)return{type:"del",raw:t[0],text:t[2],tokens:this.lexer.inlineTokens(t[2])}}autolink(e){const t=this.rules.inline.autolink.exec(e);if(t){let n,l;return t[2]==="@"?(n=ke(t[1]),l="mailto:"+n):(n=ke(t[1]),l=n),{type:"link",raw:t[0],text:n,href:l,tokens:[{type:"text",raw:n,text:n}]}}}url(e){var n;let t;if(t=this.rules.inline.url.exec(e)){let l,s;if(t[2]==="@")l=ke(t[0]),s="mailto:"+l;else{let r;do r=t[0],t[0]=((n=this.rules.inline._backpedal.exec(t[0]))==null?void 0:n[0])??"";while(r!==t[0]);l=ke(t[0]),t[1]==="www."?s="http://"+t[0]:s=t[0]}return{type:"link",raw:t[0],text:l,href:s,tokens:[{type:"text",raw:l,text:l}]}}}inlineText(e){const t=this.rules.inline.text.exec(e);if(t){let n;return this.lexer.state.inRawBlock?n=t[0]:n=ke(t[0]),{type:"text",raw:t[0],text:n}}}}const Bn=/^(?:[ \t]*(?:\n|$))+/,Pn=/^((?: {4}| {0,3}\t)[^\n]+(?:\n(?:[ \t]*(?:\n|$))*)?)+/,jn=/^ {0,3}(`{3,}(?=[^`\n]*(?:\n|$))|~{3,})([^\n]*)(?:\n|$)(?:|([\s\S]*?)(?:\n|$))(?: {0,3}\1[~`]* *(?=\n|$)|$)/,Ge=/^ {0,3}((?:-[\t ]*){3,}|(?:_[ \t]*){3,}|(?:\*[ \t]*){3,})(?:\n+|$)/,Zn=/^ {0,3}(#{1,6})(?=\s|$)(.*)(?:\n+|$)/,Vt=/(?:[*+-]|\d{1,9}[.)])/,Wt=B(/^(?!bull |blockCode|fences|blockquote|heading|html)((?:.|\n(?!\s*?\n|bull |blockCode|fences|blockquote|heading|html))+?)\n {0,3}(=+|-+) *(?:\n+|$)/).replace(/bull/g,Vt).replace(/blockCode/g,/(?: {4}| {0,3}\t)/).replace(/fences/g,/ {0,3}(?:`{3,}|~{3,})/).replace(/blockquote/g,/ {0,3}>/).replace(/heading/g,/ {0,3}#{1,6}/).replace(/html/g,/ {0,3}<[^\n>]+>\n/).getRegex(),dt=/^([^\n]+(?:\n(?!hr|heading|lheading|blockquote|fences|list|html|table| +\n)[^\n]+)*)/,Qn=/^[^\n]+/,ht=/(?!\s*\])(?:\\.|[^\[\]\\])+/,Hn=B(/^ {0,3}\[(label)\]: *(?:\n[ \t]*)?([^<\s][^\s]*|<.*?>)(?:(?: +(?:\n[ \t]*)?| *\n[ \t]*)(title))? *(?:\n+|$)/).replace("label",ht).replace("title",/(?:"(?:\\"?|[^"\\])*"|'[^'\n]*(?:\n[^'\n]+)*\n?'|\([^()]*\))/).getRegex(),Un=B(/^( {0,3}bull)([ \t][^\n]+?)?(?:\n|$)/).replace(/bull/g,Vt).getRegex(),lt="address|article|aside|base|basefont|blockquote|body|caption|center|col|colgroup|dd|details|dialog|dir|div|dl|dt|fieldset|figcaption|figure|footer|form|frame|frameset|h[1-6]|head|header|hr|html|iframe|legend|li|link|main|menu|menuitem|meta|nav|noframes|ol|optgroup|option|p|param|search|section|summary|table|tbody|td|tfoot|th|thead|title|tr|track|ul",gt=/<!--(?:-?>|[\s\S]*?(?:-->|$))/,Fn=B("^ {0,3}(?:<(script|pre|style|textarea)[\\s>][\\s\\S]*?(?:</\\1>[^\\n]*\\n+|$)|comment[^\\n]*(\\n+|$)|<\\?[\\s\\S]*?(?:\\?>\\n*|$)|<![A-Z][\\s\\S]*?(?:>\\n*|$)|<!\\[CDATA\\[[\\s\\S]*?(?:\\]\\]>\\n*|$)|</?(tag)(?: +|\\n|/?>)[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|<(?!script|pre|style|textarea)([a-z][\\w-]*)(?:attribute)*? */?>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|</(?!script|pre|style|textarea)[a-z][\\w-]*\\s*>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$))","i").replace("comment",gt).replace("tag",lt).replace("attribute",/ +[a-zA-Z:_][\w.:-]*(?: *= *"[^"\n]*"| *= *'[^'\n]*'| *= *[^\s"'=<>`]+)?/).getRegex(),Gt=B(dt).replace("hr",Ge).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("|lheading","").replace("|table","").replace("blockquote"," {0,3}>").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",lt).getRegex(),Vn=B(/^( {0,3}> ?(paragraph|[^\n]*)(?:\n|$))+/).replace("paragraph",Gt).getRegex(),mt={blockquote:Vn,code:Pn,def:Hn,fences:jn,heading:Zn,hr:Ge,html:Fn,lheading:Wt,list:Un,newline:Bn,paragraph:Gt,table:Fe,text:Qn},Et=B("^ *([^\\n ].*)\\n {0,3}((?:\\| *)?:?-+:? *(?:\\| *:?-+:? *)*(?:\\| *)?)(?:\\n((?:(?! *\\n|hr|heading|blockquote|code|fences|list|html).*(?:\\n|$))*)\\n*|$)").replace("hr",Ge).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("blockquote"," {0,3}>").replace("code","(?: {4}| {0,3}	)[^\\n]").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",lt).getRegex(),Wn={...mt,table:Et,paragraph:B(dt).replace("hr",Ge).replace("heading"," {0,3}#{1,6}(?:\\s|$)").replace("|lheading","").replace("table",Et).replace("blockquote"," {0,3}>").replace("fences"," {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~{3,})[^\\n]*\\n").replace("list"," {0,3}(?:[*+-]|1[.)]) ").replace("html","</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag",lt).getRegex()},Gn={...mt,html:B(`^ *(?:comment *(?:\\n|\\s*$)|<(tag)[\\s\\S]+?</\\1> *(?:\\n{2,}|\\s*$)|<tag(?:"[^"]*"|'[^']*'|\\s[^'"/>\\s]*)*?/?> *(?:\\n{2,}|\\s*$))`).replace("comment",gt).replace(/tag/g,"(?!(?:a|em|strong|small|s|cite|q|dfn|abbr|data|time|code|var|samp|kbd|sub|sup|i|b|u|mark|ruby|rt|rp|bdi|bdo|span|br|wbr|ins|del|img)\\b)\\w+(?!:|[^\\w\\s@]*@)\\b").getRegex(),def:/^ *\[([^\]]+)\]: *<?([^\s>]+)>?(?: +(["(][^\n]+[")]))? *(?:\n+|$)/,heading:/^(#{1,6})(.*)(?:\n+|$)/,fences:Fe,lheading:/^(.+?)\n {0,3}(=+|-+) *(?:\n+|$)/,paragraph:B(dt).replace("hr",Ge).replace("heading",` *#{1,6} *[^
]`).replace("lheading",Wt).replace("|table","").replace("blockquote"," {0,3}>").replace("|fences","").replace("|list","").replace("|html","").replace("|tag","").getRegex()},Jt=/^\\([!"#$%&'()*+,\-./:;<=>?@\[\]\\^_`{|}~])/,Jn=/^(`+)([^`]|[^`][\s\S]*?[^`])\1(?!`)/,Xt=/^( {2,}|\\)\n(?!\s*$)/,Xn=/^(`+|[^`])(?:(?= {2,}\n)|[\s\S]*?(?:(?=[\\<!\[`*_]|\b_|$)|[^ ](?= {2,}\n)))/,Je="\\p{P}\\p{S}",Yn=B(/^((?![*_])[\spunctuation])/,"u").replace(/punctuation/g,Je).getRegex(),Kn=/\[[^[\]]*?\]\((?:\\.|[^\\\(\)]|\((?:\\.|[^\\\(\)])*\))*\)|`[^`]*?`|<[^<>]*?>/g,el=B(/^(?:\*+(?:((?!\*)[punct])|[^\s*]))|^_+(?:((?!_)[punct])|([^\s_]))/,"u").replace(/punct/g,Je).getRegex(),tl=B("^[^_*]*?__[^_*]*?\\*[^_*]*?(?=__)|[^*]+(?=[^*])|(?!\\*)[punct](\\*+)(?=[\\s]|$)|[^punct\\s](\\*+)(?!\\*)(?=[punct\\s]|$)|(?!\\*)[punct\\s](\\*+)(?=[^punct\\s])|[\\s](\\*+)(?!\\*)(?=[punct])|(?!\\*)[punct](\\*+)(?!\\*)(?=[punct])|[^punct\\s](\\*+)(?=[^punct\\s])","gu").replace(/punct/g,Je).getRegex(),nl=B("^[^_*]*?\\*\\*[^_*]*?_[^_*]*?(?=\\*\\*)|[^_]+(?=[^_])|(?!_)[punct](_+)(?=[\\s]|$)|[^punct\\s](_+)(?!_)(?=[punct\\s]|$)|(?!_)[punct\\s](_+)(?=[^punct\\s])|[\\s](_+)(?!_)(?=[punct])|(?!_)[punct](_+)(?!_)(?=[punct])","gu").replace(/punct/g,Je).getRegex(),ll=B(/\\([punct])/,"gu").replace(/punct/g,Je).getRegex(),sl=B(/^<(scheme:[^\s\x00-\x1f<>]*|email)>/).replace("scheme",/[a-zA-Z][a-zA-Z0-9+.-]{1,31}/).replace("email",/[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+(@)[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+(?![-_])/).getRegex(),rl=B(gt).replace("(?:-->|$)","-->").getRegex(),ol=B("^comment|^</[a-zA-Z][\\w:-]*\\s*>|^<[a-zA-Z][\\w-]*(?:attribute)*?\\s*/?>|^<\\?[\\s\\S]*?\\?>|^<![a-zA-Z]+\\s[\\s\\S]*?>|^<!\\[CDATA\\[[\\s\\S]*?\\]\\]>").replace("comment",rl).replace("attribute",/\s+[a-zA-Z:_][\w.:-]*(?:\s*=\s*"[^"]*"|\s*=\s*'[^']*'|\s*=\s*[^\s"'=<>`]+)?/).getRegex(),tt=/(?:\[(?:\\.|[^\[\]\\])*\]|\\.|`[^`]*`|[^\[\]\\`])*?/,il=B(/^!?\[(label)\]\(\s*(href)(?:\s+(title))?\s*\)/).replace("label",tt).replace("href",/<(?:\\.|[^\n<>\\])+>|[^\s\x00-\x1f]*/).replace("title",/"(?:\\"?|[^"\\])*"|'(?:\\'?|[^'\\])*'|\((?:\\\)?|[^)\\])*\)/).getRegex(),Yt=B(/^!?\[(label)\]\[(ref)\]/).replace("label",tt).replace("ref",ht).getRegex(),Kt=B(/^!?\[(ref)\](?:\[\])?/).replace("ref",ht).getRegex(),al=B("reflink|nolink(?!\\()","g").replace("reflink",Yt).replace("nolink",Kt).getRegex(),bt={_backpedal:Fe,anyPunctuation:ll,autolink:sl,blockSkip:Kn,br:Xt,code:Jn,del:Fe,emStrongLDelim:el,emStrongRDelimAst:tl,emStrongRDelimUnd:nl,escape:Jt,link:il,nolink:Kt,punctuation:Yn,reflink:Yt,reflinkSearch:al,tag:ol,text:Xn,url:Fe},cl={...bt,link:B(/^!?\[(label)\]\((.*?)\)/).replace("label",tt).getRegex(),reflink:B(/^!?\[(label)\]\s*\[([^\]]*)\]/).replace("label",tt).getRegex()},ct={...bt,escape:B(Jt).replace("])","~|])").getRegex(),url:B(/^((?:ftp|https?):\/\/|www\.)(?:[a-zA-Z0-9\-]+\.?)+[^\s<]*|^email/,"i").replace("email",/[A-Za-z0-9._+-]+(@)[a-zA-Z0-9-_]+(?:\.[a-zA-Z0-9-_]*[a-zA-Z0-9])+(?![-_])/).getRegex(),_backpedal:/(?:[^?!.,:;*_'"~()&]+|\([^)]*\)|&(?![a-zA-Z0-9]+;$)|[?!.,:;*_'"~)]+(?!$))+/,del:/^(~~?)(?=[^\s~])((?:\\.|[^\\])*?(?:\\.|[^\s~\\]))\1(?=[^~]|$)/,text:/^([`~]+|[^`~])(?:(?= {2,}\n)|(?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)|[\s\S]*?(?:(?=[\\<!\[`*~_]|\b_|https?:\/\/|ftp:\/\/|www\.|$)|[^ ](?= {2,}\n)|[^a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-](?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)))/},ul={...ct,br:B(Xt).replace("{2,}","*").getRegex(),text:B(ct.text).replace("\\b_","\\b_| {2,}\\n").replace(/\{2,\}/g,"*").getRegex()},Ke={normal:mt,gfm:Wn,pedantic:Gn},Ue={normal:bt,gfm:ct,breaks:ul,pedantic:cl};class _e{constructor(e){P(this,"tokens");P(this,"options");P(this,"state");P(this,"tokenizer");P(this,"inlineQueue");this.tokens=[],this.tokens.links=Object.create(null),this.options=e||Me,this.options.tokenizer=this.options.tokenizer||new et,this.tokenizer=this.options.tokenizer,this.tokenizer.options=this.options,this.tokenizer.lexer=this,this.inlineQueue=[],this.state={inLink:!1,inRawBlock:!1,top:!0};const t={block:Ke.normal,inline:Ue.normal};this.options.pedantic?(t.block=Ke.pedantic,t.inline=Ue.pedantic):this.options.gfm&&(t.block=Ke.gfm,this.options.breaks?t.inline=Ue.breaks:t.inline=Ue.gfm),this.tokenizer.rules=t}static get rules(){return{block:Ke,inline:Ue}}static lex(e,t){return new _e(t).lex(e)}static lexInline(e,t){return new _e(t).inlineTokens(e)}lex(e){e=e.replace(/\r\n|\r/g,`
`),this.blockTokens(e,this.tokens);for(let t=0;t<this.inlineQueue.length;t++){const n=this.inlineQueue[t];this.inlineTokens(n.src,n.tokens)}return this.inlineQueue=[],this.tokens}blockTokens(e,t=[],n=!1){this.options.pedantic&&(e=e.replace(/\t/g,"    ").replace(/^ +$/gm,""));let l,s,r;for(;e;)if(!(this.options.extensions&&this.options.extensions.block&&this.options.extensions.block.some(i=>(l=i.call({lexer:this},e,t))?(e=e.substring(l.raw.length),t.push(l),!0):!1))){if(l=this.tokenizer.space(e)){e=e.substring(l.raw.length),l.raw.length===1&&t.length>0?t[t.length-1].raw+=`
`:t.push(l);continue}if(l=this.tokenizer.code(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&(s.type==="paragraph"||s.type==="text")?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l);continue}if(l=this.tokenizer.fences(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.heading(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.hr(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.blockquote(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.list(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.html(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.def(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&(s.type==="paragraph"||s.type==="text")?(s.raw+=`
`+l.raw,s.text+=`
`+l.raw,this.inlineQueue[this.inlineQueue.length-1].src=s.text):this.tokens.links[l.tag]||(this.tokens.links[l.tag]={href:l.href,title:l.title});continue}if(l=this.tokenizer.table(e)){e=e.substring(l.raw.length),t.push(l);continue}if(l=this.tokenizer.lheading(e)){e=e.substring(l.raw.length),t.push(l);continue}if(r=e,this.options.extensions&&this.options.extensions.startBlock){let i=1/0;const a=e.slice(1);let f;this.options.extensions.startBlock.forEach(u=>{f=u.call({lexer:this},a),typeof f=="number"&&f>=0&&(i=Math.min(i,f))}),i<1/0&&i>=0&&(r=e.substring(0,i+1))}if(this.state.top&&(l=this.tokenizer.paragraph(r))){s=t[t.length-1],n&&(s==null?void 0:s.type)==="paragraph"?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue.pop(),this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l),n=r.length!==e.length,e=e.substring(l.raw.length);continue}if(l=this.tokenizer.text(e)){e=e.substring(l.raw.length),s=t[t.length-1],s&&s.type==="text"?(s.raw+=`
`+l.raw,s.text+=`
`+l.text,this.inlineQueue.pop(),this.inlineQueue[this.inlineQueue.length-1].src=s.text):t.push(l);continue}if(e){const i="Infinite loop on byte: "+e.charCodeAt(0);if(this.options.silent){console.error(i);break}else throw new Error(i)}}return this.state.top=!0,t}inline(e,t=[]){return this.inlineQueue.push({src:e,tokens:t}),t}inlineTokens(e,t=[]){let n,l,s,r=e,i,a,f;if(this.tokens.links){const u=Object.keys(this.tokens.links);if(u.length>0)for(;(i=this.tokenizer.rules.inline.reflinkSearch.exec(r))!=null;)u.includes(i[0].slice(i[0].lastIndexOf("[")+1,-1))&&(r=r.slice(0,i.index)+"["+"a".repeat(i[0].length-2)+"]"+r.slice(this.tokenizer.rules.inline.reflinkSearch.lastIndex))}for(;(i=this.tokenizer.rules.inline.blockSkip.exec(r))!=null;)r=r.slice(0,i.index)+"["+"a".repeat(i[0].length-2)+"]"+r.slice(this.tokenizer.rules.inline.blockSkip.lastIndex);for(;(i=this.tokenizer.rules.inline.anyPunctuation.exec(r))!=null;)r=r.slice(0,i.index)+"++"+r.slice(this.tokenizer.rules.inline.anyPunctuation.lastIndex);for(;e;)if(a||(f=""),a=!1,!(this.options.extensions&&this.options.extensions.inline&&this.options.extensions.inline.some(u=>(n=u.call({lexer:this},e,t))?(e=e.substring(n.raw.length),t.push(n),!0):!1))){if(n=this.tokenizer.escape(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.tag(e)){e=e.substring(n.raw.length),l=t[t.length-1],l&&n.type==="text"&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(n=this.tokenizer.link(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.reflink(e,this.tokens.links)){e=e.substring(n.raw.length),l=t[t.length-1],l&&n.type==="text"&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(n=this.tokenizer.emStrong(e,r,f)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.codespan(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.br(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.del(e)){e=e.substring(n.raw.length),t.push(n);continue}if(n=this.tokenizer.autolink(e)){e=e.substring(n.raw.length),t.push(n);continue}if(!this.state.inLink&&(n=this.tokenizer.url(e))){e=e.substring(n.raw.length),t.push(n);continue}if(s=e,this.options.extensions&&this.options.extensions.startInline){let u=1/0;const d=e.slice(1);let c;this.options.extensions.startInline.forEach(_=>{c=_.call({lexer:this},d),typeof c=="number"&&c>=0&&(u=Math.min(u,c))}),u<1/0&&u>=0&&(s=e.substring(0,u+1))}if(n=this.tokenizer.inlineText(s)){e=e.substring(n.raw.length),n.raw.slice(-1)!=="_"&&(f=n.raw.slice(-1)),a=!0,l=t[t.length-1],l&&l.type==="text"?(l.raw+=n.raw,l.text+=n.text):t.push(n);continue}if(e){const u="Infinite loop on byte: "+e.charCodeAt(0);if(this.options.silent){console.error(u);break}else throw new Error(u)}}return t}}class nt{constructor(e){P(this,"options");P(this,"parser");this.options=e||Me}space(e){return""}code({text:e,lang:t,escaped:n}){var r;const l=(r=(t||"").match(/^\S*/))==null?void 0:r[0],s=e.replace(/\n$/,"")+`
`;return l?'<pre><code class="language-'+ke(l)+'">'+(n?s:ke(s,!0))+`</code></pre>
`:"<pre><code>"+(n?s:ke(s,!0))+`</code></pre>
`}blockquote({tokens:e}){return`<blockquote>
${this.parser.parse(e)}</blockquote>
`}html({text:e}){return e}heading({tokens:e,depth:t}){return`<h${t}>${this.parser.parseInline(e)}</h${t}>
`}hr(e){return`<hr>
`}list(e){const t=e.ordered,n=e.start;let l="";for(let i=0;i<e.items.length;i++){const a=e.items[i];l+=this.listitem(a)}const s=t?"ol":"ul",r=t&&n!==1?' start="'+n+'"':"";return"<"+s+r+`>
`+l+"</"+s+`>
`}listitem(e){let t="";if(e.task){const n=this.checkbox({checked:!!e.checked});e.loose?e.tokens.length>0&&e.tokens[0].type==="paragraph"?(e.tokens[0].text=n+" "+e.tokens[0].text,e.tokens[0].tokens&&e.tokens[0].tokens.length>0&&e.tokens[0].tokens[0].type==="text"&&(e.tokens[0].tokens[0].text=n+" "+e.tokens[0].tokens[0].text)):e.tokens.unshift({type:"text",raw:n+" ",text:n+" "}):t+=n+" "}return t+=this.parser.parse(e.tokens,!!e.loose),`<li>${t}</li>
`}checkbox({checked:e}){return"<input "+(e?'checked="" ':"")+'disabled="" type="checkbox">'}paragraph({tokens:e}){return`<p>${this.parser.parseInline(e)}</p>
`}table(e){let t="",n="";for(let s=0;s<e.header.length;s++)n+=this.tablecell(e.header[s]);t+=this.tablerow({text:n});let l="";for(let s=0;s<e.rows.length;s++){const r=e.rows[s];n="";for(let i=0;i<r.length;i++)n+=this.tablecell(r[i]);l+=this.tablerow({text:n})}return l&&(l=`<tbody>${l}</tbody>`),`<table>
<thead>
`+t+`</thead>
`+l+`</table>
`}tablerow({text:e}){return`<tr>
${e}</tr>
`}tablecell(e){const t=this.parser.parseInline(e.tokens),n=e.header?"th":"td";return(e.align?`<${n} align="${e.align}">`:`<${n}>`)+t+`</${n}>
`}strong({tokens:e}){return`<strong>${this.parser.parseInline(e)}</strong>`}em({tokens:e}){return`<em>${this.parser.parseInline(e)}</em>`}codespan({text:e}){return`<code>${e}</code>`}br(e){return"<br>"}del({tokens:e}){return`<del>${this.parser.parseInline(e)}</del>`}link({href:e,title:t,tokens:n}){const l=this.parser.parseInline(n),s=yt(e);if(s===null)return l;e=s;let r='<a href="'+e+'"';return t&&(r+=' title="'+t+'"'),r+=">"+l+"</a>",r}image({href:e,title:t,text:n}){const l=yt(e);if(l===null)return n;e=l;let s=`<img src="${e}" alt="${n}"`;return t&&(s+=` title="${t}"`),s+=">",s}text(e){return"tokens"in e&&e.tokens?this.parser.parseInline(e.tokens):e.text}}class kt{strong({text:e}){return e}em({text:e}){return e}codespan({text:e}){return e}del({text:e}){return e}html({text:e}){return e}text({text:e}){return e}link({text:e}){return""+e}image({text:e}){return""+e}br(){return""}}class we{constructor(e){P(this,"options");P(this,"renderer");P(this,"textRenderer");this.options=e||Me,this.options.renderer=this.options.renderer||new nt,this.renderer=this.options.renderer,this.renderer.options=this.options,this.renderer.parser=this,this.textRenderer=new kt}static parse(e,t){return new we(t).parse(e)}static parseInline(e,t){return new we(t).parseInline(e)}parse(e,t=!0){let n="";for(let l=0;l<e.length;l++){const s=e[l];if(this.options.extensions&&this.options.extensions.renderers&&this.options.extensions.renderers[s.type]){const i=s,a=this.options.extensions.renderers[i.type].call({parser:this},i);if(a!==!1||!["space","hr","heading","code","table","blockquote","list","html","paragraph","text"].includes(i.type)){n+=a||"";continue}}const r=s;switch(r.type){case"space":{n+=this.renderer.space(r);continue}case"hr":{n+=this.renderer.hr(r);continue}case"heading":{n+=this.renderer.heading(r);continue}case"code":{n+=this.renderer.code(r);continue}case"table":{n+=this.renderer.table(r);continue}case"blockquote":{n+=this.renderer.blockquote(r);continue}case"list":{n+=this.renderer.list(r);continue}case"html":{n+=this.renderer.html(r);continue}case"paragraph":{n+=this.renderer.paragraph(r);continue}case"text":{let i=r,a=this.renderer.text(i);for(;l+1<e.length&&e[l+1].type==="text";)i=e[++l],a+=`
`+this.renderer.text(i);t?n+=this.renderer.paragraph({type:"paragraph",raw:a,text:a,tokens:[{type:"text",raw:a,text:a}]}):n+=a;continue}default:{const i='Token with "'+r.type+'" type was not found.';if(this.options.silent)return console.error(i),"";throw new Error(i)}}}return n}parseInline(e,t){t=t||this.renderer;let n="";for(let l=0;l<e.length;l++){const s=e[l];if(this.options.extensions&&this.options.extensions.renderers&&this.options.extensions.renderers[s.type]){const i=this.options.extensions.renderers[s.type].call({parser:this},s);if(i!==!1||!["escape","html","link","image","strong","em","codespan","br","del","text"].includes(s.type)){n+=i||"";continue}}const r=s;switch(r.type){case"escape":{n+=t.text(r);break}case"html":{n+=t.html(r);break}case"link":{n+=t.link(r);break}case"image":{n+=t.image(r);break}case"strong":{n+=t.strong(r);break}case"em":{n+=t.em(r);break}case"codespan":{n+=t.codespan(r);break}case"br":{n+=t.br(r);break}case"del":{n+=t.del(r);break}case"text":{n+=t.text(r);break}default:{const i='Token with "'+r.type+'" type was not found.';if(this.options.silent)return console.error(i),"";throw new Error(i)}}}return n}}class Ve{constructor(e){P(this,"options");P(this,"block");this.options=e||Me}preprocess(e){return e}postprocess(e){return e}processAllTokens(e){return e}provideLexer(){return this.block?_e.lex:_e.lexInline}provideParser(){return this.block?we.parse:we.parseInline}}P(Ve,"passThroughHooks",new Set(["preprocess","postprocess","processAllTokens"]));class fl{constructor(...e){P(this,"defaults",pt());P(this,"options",this.setOptions);P(this,"parse",this.parseMarkdown(!0));P(this,"parseInline",this.parseMarkdown(!1));P(this,"Parser",we);P(this,"Renderer",nt);P(this,"TextRenderer",kt);P(this,"Lexer",_e);P(this,"Tokenizer",et);P(this,"Hooks",Ve);this.use(...e)}walkTokens(e,t){var l,s;let n=[];for(const r of e)switch(n=n.concat(t.call(this,r)),r.type){case"table":{const i=r;for(const a of i.header)n=n.concat(this.walkTokens(a.tokens,t));for(const a of i.rows)for(const f of a)n=n.concat(this.walkTokens(f.tokens,t));break}case"list":{const i=r;n=n.concat(this.walkTokens(i.items,t));break}default:{const i=r;(s=(l=this.defaults.extensions)==null?void 0:l.childTokens)!=null&&s[i.type]?this.defaults.extensions.childTokens[i.type].forEach(a=>{const f=i[a].flat(1/0);n=n.concat(this.walkTokens(f,t))}):i.tokens&&(n=n.concat(this.walkTokens(i.tokens,t)))}}return n}use(...e){const t=this.defaults.extensions||{renderers:{},childTokens:{}};return e.forEach(n=>{const l={...n};if(l.async=this.defaults.async||l.async||!1,n.extensions&&(n.extensions.forEach(s=>{if(!s.name)throw new Error("extension name required");if("renderer"in s){const r=t.renderers[s.name];r?t.renderers[s.name]=function(...i){let a=s.renderer.apply(this,i);return a===!1&&(a=r.apply(this,i)),a}:t.renderers[s.name]=s.renderer}if("tokenizer"in s){if(!s.level||s.level!=="block"&&s.level!=="inline")throw new Error("extension level must be 'block' or 'inline'");const r=t[s.level];r?r.unshift(s.tokenizer):t[s.level]=[s.tokenizer],s.start&&(s.level==="block"?t.startBlock?t.startBlock.push(s.start):t.startBlock=[s.start]:s.level==="inline"&&(t.startInline?t.startInline.push(s.start):t.startInline=[s.start]))}"childTokens"in s&&s.childTokens&&(t.childTokens[s.name]=s.childTokens)}),l.extensions=t),n.renderer){const s=this.defaults.renderer||new nt(this.defaults);for(const r in n.renderer){if(!(r in s))throw new Error(`renderer '${r}' does not exist`);if(["options","parser"].includes(r))continue;const i=r,a=n.renderer[i],f=s[i];s[i]=(...u)=>{let d=a.apply(s,u);return d===!1&&(d=f.apply(s,u)),d||""}}l.renderer=s}if(n.tokenizer){const s=this.defaults.tokenizer||new et(this.defaults);for(const r in n.tokenizer){if(!(r in s))throw new Error(`tokenizer '${r}' does not exist`);if(["options","rules","lexer"].includes(r))continue;const i=r,a=n.tokenizer[i],f=s[i];s[i]=(...u)=>{let d=a.apply(s,u);return d===!1&&(d=f.apply(s,u)),d}}l.tokenizer=s}if(n.hooks){const s=this.defaults.hooks||new Ve;for(const r in n.hooks){if(!(r in s))throw new Error(`hook '${r}' does not exist`);if(["options","block"].includes(r))continue;const i=r,a=n.hooks[i],f=s[i];Ve.passThroughHooks.has(r)?s[i]=u=>{if(this.defaults.async)return Promise.resolve(a.call(s,u)).then(c=>f.call(s,c));const d=a.call(s,u);return f.call(s,d)}:s[i]=(...u)=>{let d=a.apply(s,u);return d===!1&&(d=f.apply(s,u)),d}}l.hooks=s}if(n.walkTokens){const s=this.defaults.walkTokens,r=n.walkTokens;l.walkTokens=function(i){let a=[];return a.push(r.call(this,i)),s&&(a=a.concat(s.call(this,i))),a}}this.defaults={...this.defaults,...l}}),this}setOptions(e){return this.defaults={...this.defaults,...e},this}lexer(e,t){return _e.lex(e,t??this.defaults)}parser(e,t){return we.parse(e,t??this.defaults)}parseMarkdown(e){return(n,l)=>{const s={...l},r={...this.defaults,...s},i=this.onError(!!r.silent,!!r.async);if(this.defaults.async===!0&&s.async===!1)return i(new Error("marked(): The async option was set to true by an extension. Remove async: false from the parse options object to return a Promise."));if(typeof n>"u"||n===null)return i(new Error("marked(): input parameter is undefined or null"));if(typeof n!="string")return i(new Error("marked(): input parameter is of type "+Object.prototype.toString.call(n)+", string expected"));r.hooks&&(r.hooks.options=r,r.hooks.block=e);const a=r.hooks?r.hooks.provideLexer():e?_e.lex:_e.lexInline,f=r.hooks?r.hooks.provideParser():e?we.parse:we.parseInline;if(r.async)return Promise.resolve(r.hooks?r.hooks.preprocess(n):n).then(u=>a(u,r)).then(u=>r.hooks?r.hooks.processAllTokens(u):u).then(u=>r.walkTokens?Promise.all(this.walkTokens(u,r.walkTokens)).then(()=>u):u).then(u=>f(u,r)).then(u=>r.hooks?r.hooks.postprocess(u):u).catch(i);try{r.hooks&&(n=r.hooks.preprocess(n));let u=a(n,r);r.hooks&&(u=r.hooks.processAllTokens(u)),r.walkTokens&&this.walkTokens(u,r.walkTokens);let d=f(u,r);return r.hooks&&(d=r.hooks.postprocess(d)),d}catch(u){return i(u)}}}onError(e,t){return n=>{if(n.message+=`
Please report this to https://github.com/markedjs/marked.`,e){const l="<p>An error occurred:</p><pre>"+ke(n.message+"",!0)+"</pre>";return t?Promise.resolve(l):l}if(t)return Promise.reject(n);throw n}}}const qe=new fl;function N(o,e){return qe.parse(o,e)}N.options=N.setOptions=function(o){return qe.setOptions(o),N.defaults=qe.defaults,Ht(N.defaults),N};N.getDefaults=pt;N.defaults=Me;N.use=function(...o){return qe.use(...o),N.defaults=qe.defaults,Ht(N.defaults),N};N.walkTokens=function(o,e){return qe.walkTokens(o,e)};N.parseInline=qe.parseInline;N.Parser=we;N.parser=we.parse;N.Renderer=nt;N.TextRenderer=kt;N.Lexer=_e;N.lexer=_e.lex;N.Tokenizer=et;N.Hooks=Ve;N.parse=N;N.options;N.setOptions;N.use;N.walkTokens;N.parseInline;we.parse;_e.lex;function vt(o,e,t){const n=o.slice();return n[67]=e[t],n[68]=e,n[69]=t,n}function Rt(o,e,t){const n=o.slice();return n[70]=e[t],n}function Ct(o,e,t){const n=o.slice();return n[73]=e[t],n}function St(o,e,t){const n=o.slice();return n[76]=e[t],n}function pl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"play",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function dl(o){let e;return{c(){e=y("span"),g(e,"class","loading loading-spinner")},m(t,n){$(t,e,n)},i:Ie,o:Ie,d(t){t&&I(e)}}}function $t(o){let e,t,n,l,s,r;const i=[gl,hl],a=[];function f(d,c){return d[67].isEditing?0:1}t=f(o),n=a[t]=i[t](o);function u(){return o[35](o[67])}return{c(){e=y("button"),n.c(),g(e,"class","btn btn-ghost btn-xs"),g(e,"title","Toggle edit mode"),he(e,"active",o[67].isEditing)},m(d,c){$(d,e,c),a[t].m(e,null),l=!0,s||(r=ee(e,"click",u),s=!0)},p(d,c){o=d;let _=t;t=f(o),t!==_&&(oe(),v(a[_],1,1,()=>{a[_]=null}),ie(),n=a[t],n||(n=a[t]=i[t](o),n.c()),T(n,1),n.m(e,null)),(!l||c[0]&32)&&he(e,"active",o[67].isEditing)},i(d){l||(T(n),l=!0)},o(d){v(n),l=!1},d(d){d&&I(e),a[t].d(),s=!1,r()}}}function hl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"pencil",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function gl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"eye",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function ml(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"eye-slash",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function bl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"eye",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function It(o){let e,t,n,l,s,r,i,a,f,u,d,c,_,h,m,w,R,A,x,j,te;const q=[_l,kl],W=[];function xe(C,Q){return C[67].isExecuting?0:1}t=xe(o),n=W[t]=q[t](o);function ne(){return o[37](o[67])}const le=[xl,wl],U=[];function me(C,Q){return C[67].showOutput?0:1}i=me(o),a=U[i]=le[i](o);function be(){return o[38](o[67])}c=new ae({props:{type:"outline",icon:"chevron-up",path:"icons"}});function se(){return o[39](o[67])}w=new ae({props:{type:"outline",icon:"chevron-down",path:"icons"}});function ue(){return o[40](o[67])}function ye(){o[41].call(A,o[68],o[69])}return{c(){e=y("button"),n.c(),s=D(),r=y("button"),a.c(),f=D(),u=y("button"),d=y("span"),G(c.$$.fragment),_=D(),h=y("button"),m=y("span"),G(w.$$.fragment),R=D(),A=y("input"),g(e,"class","btn btn-ghost btn-xs"),e.disabled=l=o[67].isExecuting,g(e,"title","Run this cell"),he(e,"executing",o[67].isExecuting),g(r,"class","btn btn-ghost btn-xs"),g(r,"title","Toggle output visibility"),he(r,"active",o[67].showOutput),g(d,"class","w-5 h-5"),g(u,"class","btn btn-ghost btn-xs"),g(u,"title","Run all cells above"),g(m,"class","w-5 h-5"),g(h,"class","btn btn-ghost btn-xs"),g(h,"title","Run all cells below"),g(A,"type","text"),g(A,"placeholder","Variable Name"),g(A,"class","input input-xs")},m(C,Q){$(C,e,Q),W[t].m(e,null),$(C,s,Q),$(C,r,Q),U[i].m(r,null),$(C,f,Q),$(C,u,Q),k(u,d),J(c,d,null),$(C,_,Q),$(C,h,Q),k(h,m),J(w,m,null),$(C,R,Q),$(C,A,Q),wt(A,o[67].name),x=!0,j||(te=[ee(e,"click",ne),ee(r,"click",be),ee(u,"click",se),ee(h,"click",ue),ee(A,"input",ye)],j=!0)},p(C,Q){o=C;let Y=t;t=xe(o),t!==Y&&(oe(),v(W[Y],1,1,()=>{W[Y]=null}),ie(),n=W[t],n||(n=W[t]=q[t](o),n.c()),T(n,1),n.m(e,null)),(!x||Q[0]&32&&l!==(l=o[67].isExecuting))&&(e.disabled=l),(!x||Q[0]&32)&&he(e,"executing",o[67].isExecuting);let Te=i;i=me(o),i!==Te&&(oe(),v(U[Te],1,1,()=>{U[Te]=null}),ie(),a=U[i],a||(a=U[i]=le[i](o),a.c()),T(a,1),a.m(r,null)),(!x||Q[0]&32)&&he(r,"active",o[67].showOutput),Q[0]&32&&A.value!==o[67].name&&wt(A,o[67].name)},i(C){x||(T(n),T(a),T(c.$$.fragment,C),T(w.$$.fragment,C),x=!0)},o(C){v(n),v(a),v(c.$$.fragment,C),v(w.$$.fragment,C),x=!1},d(C){C&&(I(e),I(s),I(r),I(f),I(u),I(_),I(h),I(R),I(A)),W[t].d(),U[i].d(),X(c),X(w),j=!1,ut(te)}}}function kl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"play",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function _l(o){let e;return{c(){e=y("span"),g(e,"class","loading loading-spinner")},m(t,n){$(t,e,n)},i:Ie,o:Ie,d(t){t&&I(e)}}}function wl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"chart-bar",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function xl(o){let e,t,n;return t=new ae({props:{type:"outline",icon:"chart-bar",path:"icons"}}),{c(){e=y("span"),G(t.$$.fragment),g(e,"class","w-5 h-5")},m(l,s){$(l,e,s),J(t,e,null),n=!0},i(l){n||(T(t.$$.fragment,l),n=!0)},o(l){v(t.$$.fragment,l),n=!1},d(l){l&&I(e),X(t)}}}function yl(o){let e,t,n=o[67].showCode&&At(o);return{c(){n&&n.c(),e=ft()},m(l,s){n&&n.m(l,s),$(l,e,s),t=!0},p(l,s){l[67].showCode?n?(n.p(l,s),s[0]&32&&T(n,1)):(n=At(l),n.c(),T(n,1),n.m(e.parentNode,e)):n&&(oe(),v(n,1,1,()=>{n=null}),ie())},i(l){t||(T(n),t=!0)},o(l){v(n),t=!1},d(l){l&&I(e),n&&n.d(l)}}}function Tl(o){let e,t,n,l;const s=[El,zl],r=[];function i(a,f){return a[67].isEditing&&a[67].showCode?0:!a[67].isEditing&&a[67].showCode?1:-1}return~(e=i(o))&&(t=r[e]=s[e](o)),{c(){t&&t.c(),n=ft()},m(a,f){~e&&r[e].m(a,f),$(a,n,f),l=!0},p(a,f){let u=e;e=i(a),e===u?~e&&r[e].p(a,f):(t&&(oe(),v(r[u],1,1,()=>{r[u]=null}),ie()),~e?(t=r[e],t?t.p(a,f):(t=r[e]=s[e](a),t.c()),T(t,1),t.m(n.parentNode,n)):t=null)},i(a){l||(T(t),l=!0)},o(a){v(t),l=!1},d(a){a&&I(n),~e&&r[e].d(a)}}}function At(o){let e,t,n,l;function s(...r){return o[45](o[67],...r)}return t=new Qt({props:{code:o[67].code,language:o[67].type,minHeight:30,maxHeight:700,onEditStop:s,options:{fontSize:12,wordWrap:"on",lineNumbers:"on",minimap:{enabled:!1},tabSize:4}}}),{c(){e=y("div"),n=y("div"),G(t.$$.fragment),We(n,"display","contents"),We(n,"--theme","vs-dark"),g(e,"class","cell-content")},m(r,i){$(r,e,i),k(e,n),J(t,n,null),l=!0},p(r,i){o=r;const a={};i[0]&32&&(a.code=o[67].code),i[0]&32&&(a.language=o[67].type),i[0]&32&&(a.onEditStop=s),t.$set(a)},i(r){l||(T(t.$$.fragment,r),l=!0)},o(r){v(t.$$.fragment,r),l=!1},d(r){r&&I(e),X(t)}}}function zl(o){let e,t=o[15](o[67].code)+"",n,l;function s(){return o[44](o[67])}return{c(){e=y("div"),g(e,"class","w-full p-0 markdown-rendered")},m(r,i){$(r,e,i),e.innerHTML=t,n||(l=ee(e,"dblclick",s),n=!0)},p(r,i){o=r,i[0]&32&&t!==(t=o[15](o[67].code)+"")&&(e.innerHTML=t)},i:Ie,o:Ie,d(r){r&&I(e),n=!1,l()}}}function El(o){let e,t,n,l;function s(...r){return o[43](o[67],...r)}return t=new Qt({props:{code:o[67].code,language:o[67].type,minHeight:30,maxHeight:700,onEditStop:s,options:{fontSize:12,wordWrap:"on",lineNumbers:"on",minimap:{enabled:!1},tabSize:4}}}),{c(){e=y("div"),n=y("div"),G(t.$$.fragment),We(n,"display","contents"),We(n,"--theme","vs-dark"),g(e,"class","cell-content")},m(r,i){$(r,e,i),k(e,n),J(t,n,null),l=!0},p(r,i){o=r;const a={};i[0]&32&&(a.code=o[67].code),i[0]&32&&(a.language=o[67].type),i[0]&32&&(a.onEditStop=s),t.$set(a)},i(r){l||(T(t.$$.fragment,r),l=!0)},o(r){v(t.$$.fragment,r),l=!1},d(r){r&&I(e),X(t)}}}function Lt(o){let e,t,n,l;const s=[Rl,vl],r=[];function i(a,f){return a[67].output.error?0:a[67].output.executionTime?1:-1}return~(t=i(o))&&(n=r[t]=s[t](o)),{c(){e=y("div"),n&&n.c(),g(e,"class","cell-output")},m(a,f){$(a,e,f),~t&&r[t].m(e,null),l=!0},p(a,f){let u=t;t=i(a),t===u?~t&&r[t].p(a,f):(n&&(oe(),v(r[u],1,1,()=>{r[u]=null}),ie()),~t?(n=r[t],n?n.p(a,f):(n=r[t]=s[t](a),n.c()),T(n,1),n.m(e,null)):n=null)},i(a){l||(T(n),l=!0)},o(a){v(n),l=!1},d(a){a&&I(e),~t&&r[t].d()}}}function vl(o){var h;let e,t,n,l,s=o[67].output.executionTime+"",r,i,a,f,u,d,c=((h=o[67].output)==null?void 0:h.data)&&Ot(o),_=o[67].output.data&&qt(o);return{c(){e=y("div"),t=y("span"),t.textContent="✓ Query executed successfully",n=D(),l=y("span"),r=ge(s),i=ge("ms"),a=D(),c&&c.c(),f=D(),_&&_.c(),u=ft(),g(l,"class","execution-time"),g(e,"class","alert alert-success text-sm p-1")},m(m,w){$(m,e,w),k(e,t),k(e,n),k(e,l),k(l,r),k(l,i),k(e,a),c&&c.m(e,null),$(m,f,w),_&&_.m(m,w),$(m,u,w),d=!0},p(m,w){var R;(!d||w[0]&32)&&s!==(s=m[67].output.executionTime+"")&&Pe(r,s),(R=m[67].output)!=null&&R.data?c?(c.p(m,w),w[0]&32&&T(c,1)):(c=Ot(m),c.c(),T(c,1),c.m(e,null)):c&&(oe(),v(c,1,1,()=>{c=null}),ie()),m[67].output.data?_?_.p(m,w):(_=qt(m),_.c(),_.m(u.parentNode,u)):_&&(_.d(1),_=null)},i(m){d||(T(c),d=!0)},o(m){v(c),d=!1},d(m){m&&(I(e),I(f),I(u)),c&&c.d(),_&&_.d(m)}}}function Rl(o){let e,t,n,l=o[67].output.error+"",s,r,i,a=o[67].output.executionTime+"",f,u;return{c(){e=y("div"),t=y("strong"),t.textContent="Error:",n=D(),s=ge(l),r=D(),i=y("span"),f=ge(a),u=ge("ms"),g(i,"class","execution-time"),g(e,"class","alert alert-error text-sm p-1")},m(d,c){$(d,e,c),k(e,t),k(e,n),k(e,s),k(e,r),k(e,i),k(i,f),k(i,u)},p(d,c){c[0]&32&&l!==(l=d[67].output.error+"")&&Pe(s,l),c[0]&32&&a!==(a=d[67].output.executionTime+"")&&Pe(f,a)},i:Ie,o:Ie,d(d){d&&I(e)}}}function Ot(o){let e,t,n,l,s,r;n=new ae({props:{type:"outline",icon:"arrow-down-tray",path:"icons"}});function i(){return o[46](o[67])}return{c(){e=y("button"),t=y("span"),G(n.$$.fragment),g(t,"class","w-5 h-5"),g(e,"class","btn btn-ghost btn-xs btn-export"),g(e,"title","Export as CSV")},m(a,f){$(a,e,f),k(e,t),J(n,t,null),l=!0,s||(r=ee(e,"click",i),s=!0)},p(a,f){o=a},i(a){l||(T(n.$$.fragment,a),l=!0)},o(a){v(n.$$.fragment,a),l=!1},d(a){a&&I(e),X(n),s=!1,r()}}}function qt(o){let e,t,n,l,s,r,i,a,f,u=o[67].output.data.length+"",d,c,_=$e(Object.keys(o[67].output.data[0])),h=[];for(let R=0;R<_.length;R+=1)h[R]=Mt(St(o,_,R));let m=$e(o[67].output.data.slice(0,10)),w=[];for(let R=0;R<m.length;R+=1)w[R]=Nt(Rt(o,m,R));return{c(){e=y("div"),t=y("table"),n=y("thead"),l=y("tr");for(let R=0;R<h.length;R+=1)h[R].c();s=D(),r=y("tbody");for(let R=0;R<w.length;R+=1)w[R].c();i=D(),a=y("div"),f=ge("Showing first 10 of "),d=ge(u),c=ge(" rows"),g(t,"class","data-table"),g(a,"class","table-footer"),g(e,"class","data-table-container")},m(R,A){$(R,e,A),k(e,t),k(t,n),k(n,l);for(let x=0;x<h.length;x+=1)h[x]&&h[x].m(l,null);k(t,s),k(t,r);for(let x=0;x<w.length;x+=1)w[x]&&w[x].m(r,null);k(e,i),k(e,a),k(a,f),k(a,d),k(a,c)},p(R,A){if(A[0]&32){_=$e(Object.keys(R[67].output.data[0]));let x;for(x=0;x<_.length;x+=1){const j=St(R,_,x);h[x]?h[x].p(j,A):(h[x]=Mt(j),h[x].c(),h[x].m(l,null))}for(;x<h.length;x+=1)h[x].d(1);h.length=_.length}if(A[0]&32){m=$e(R[67].output.data.slice(0,10));let x;for(x=0;x<m.length;x+=1){const j=Rt(R,m,x);w[x]?w[x].p(j,A):(w[x]=Nt(j),w[x].c(),w[x].m(r,null))}for(;x<w.length;x+=1)w[x].d(1);w.length=m.length}A[0]&32&&u!==(u=R[67].output.data.length+"")&&Pe(d,u)},d(R){R&&I(e),at(h,R),at(w,R)}}}function Mt(o){let e,t=o[76]+"",n;return{c(){e=y("th"),n=ge(t)},m(l,s){$(l,e,s),k(e,n)},p(l,s){s[0]&32&&t!==(t=l[76]+"")&&Pe(n,t)},d(l){l&&I(e)}}}function Dt(o){let e,t=o[73]+"",n;return{c(){e=y("td"),n=ge(t)},m(l,s){$(l,e,s),k(e,n)},p(l,s){s[0]&32&&t!==(t=l[73]+"")&&Pe(n,t)},d(l){l&&I(e)}}}function Nt(o){let e,t,n=$e(Object.values(o[70])),l=[];for(let s=0;s<n.length;s+=1)l[s]=Dt(Ct(o,n,s));return{c(){e=y("tr");for(let s=0;s<l.length;s+=1)l[s].c();t=D()},m(s,r){$(s,e,r);for(let i=0;i<l.length;i+=1)l[i]&&l[i].m(e,null);k(e,t)},p(s,r){if(r[0]&32){n=$e(Object.values(s[70]));let i;for(i=0;i<n.length;i+=1){const a=Ct(s,n,i);l[i]?l[i].p(a,r):(l[i]=Dt(a),l[i].c(),l[i].m(e,t))}for(;i<l.length;i+=1)l[i].d(1);l.length=n.length}},d(s){s&&I(e),at(l,s)}}}function Bt(o,e){let t,n,l,s,r,i,a,f,u,d,c,_,h,m,w,R,A,x,j,te,q,W,xe,ne,le,U,me,be,se,ue,ye,C=e[67].type==="markdown"&&$t(e);const Q=[bl,ml],Y=[];function Te(Z,re){return Z[67].showCode?0:1}r=Te(e),i=Y[r]=Q[r](e);function Re(){return e[36](e[67])}let F=e[67].type==="sql"&&It(e);c=new ae({props:{type:"outline",icon:"trash",path:"icons"}});function V(){return e[42](e[67])}const H=[Tl,yl],M=[];function E(Z,re){return Z[67].type==="markdown"?0:Z[67].type==="sql"?1:-1}~(m=E(e))&&(w=M[m]=H[m](e));let z=e[67].type==="sql"&&e[67].output&&e[67].showOutput&&Lt(e);q=new ae({props:{type:"outline",icon:"plus",path:"icons"}});function fe(){return e[47](e[69])}U=new ae({props:{type:"outline",icon:"plus",path:"icons"}});function st(){return e[48](e[69])}function rt(...Z){return e[49](e[69],...Z)}function ot(...Z){return e[50](e[69],...Z)}function it(...Z){return e[51](e[69],...Z)}return{key:o,first:null,c(){t=y("div"),n=y("div"),C&&C.c(),l=D(),s=y("button"),i.c(),a=D(),F&&F.c(),f=D(),u=y("button"),d=y("span"),G(c.$$.fragment),_=D(),h=y("div"),w&&w.c(),R=D(),z&&z.c(),A=D(),x=y("div"),j=y("button"),te=y("span"),G(q.$$.fragment),W=ge(" Code"),xe=D(),ne=y("button"),le=y("span"),G(U.$$.fragment),me=ge(" Markdown"),be=D(),g(s,"class","btn btn-ghost btn-xs"),g(s,"title","Toggle code visibility"),he(s,"active",e[67].showCode),g(d,"class","w-5 h-5"),g(u,"class","btn btn-ghost btn-xs"),g(u,"title","Delete this cell"),g(n,"class","absolute left -top-5 transform --translate-x-1/4 z-10 hidden group-hover:flex group-focus-within:flex border border-base-300 bg-base-100 px-0 py-0 space-x-1"),g(h,"class","space-y-4"),g(te,"class","w-5 h-5"),g(j,"class","btn btn-ghost btn-xs"),g(j,"title","Add SQL Cell"),g(le,"class","w-5 h-5"),g(ne,"class","btn btn-ghost btn-xs"),g(ne,"title","Add Markdown Cell"),g(x,"class","absolute left-1/2 -bottom-5 transform -translate-x-1/2 z-10 hidden group-hover:flex group-focus-within:flex border border-base-300 bg-base-100 px-0 py-0 space-x-1"),g(t,"class","relative group border border-base-200 rounded-xs p-2 pt-2 mb-4 bg-base-100"),g(t,"draggable","true"),g(t,"role","listitem"),g(t,"tabindex","0"),he(t,"dragging",e[7]===e[69]),he(t,"drag-over",e[8]===e[69]),this.first=t},m(Z,re){$(Z,t,re),k(t,n),C&&C.m(n,null),k(n,l),k(n,s),Y[r].m(s,null),k(n,a),F&&F.m(n,null),k(n,f),k(n,u),k(u,d),J(c,d,null),k(t,_),k(t,h),~m&&M[m].m(h,null),k(t,R),z&&z.m(t,null),k(t,A),k(t,x),k(x,j),k(j,te),J(q,te,null),k(j,W),k(x,xe),k(x,ne),k(ne,le),J(U,le,null),k(ne,me),k(t,be),se=!0,ue||(ye=[ee(s,"click",Re),ee(u,"click",V),ee(j,"click",fe),ee(ne,"click",st),ee(t,"dragstart",rt),ee(t,"dragover",ot),ee(t,"dragleave",e[25]),ee(t,"drop",it),ee(t,"dragend",e[27])],ue=!0)},p(Z,re){e=Z,e[67].type==="markdown"?C?(C.p(e,re),re[0]&32&&T(C,1)):(C=$t(e),C.c(),T(C,1),C.m(n,l)):C&&(oe(),v(C,1,1,()=>{C=null}),ie());let je=r;r=Te(e),r!==je&&(oe(),v(Y[je],1,1,()=>{Y[je]=null}),ie(),i=Y[r],i||(i=Y[r]=Q[r](e),i.c()),T(i,1),i.m(s,null)),(!se||re[0]&32)&&he(s,"active",e[67].showCode),e[67].type==="sql"?F?(F.p(e,re),re[0]&32&&T(F,1)):(F=It(e),F.c(),T(F,1),F.m(n,f)):F&&(oe(),v(F,1,1,()=>{F=null}),ie());let Ze=m;m=E(e),m===Ze?~m&&M[m].p(e,re):(w&&(oe(),v(M[Ze],1,1,()=>{M[Ze]=null}),ie()),~m?(w=M[m],w?w.p(e,re):(w=M[m]=H[m](e),w.c()),T(w,1),w.m(h,null)):w=null),e[67].type==="sql"&&e[67].output&&e[67].showOutput?z?(z.p(e,re),re[0]&32&&T(z,1)):(z=Lt(e),z.c(),T(z,1),z.m(t,A)):z&&(oe(),v(z,1,1,()=>{z=null}),ie()),(!se||re[0]&160)&&he(t,"dragging",e[7]===e[69]),(!se||re[0]&288)&&he(t,"drag-over",e[8]===e[69])},i(Z){se||(T(C),T(i),T(F),T(c.$$.fragment,Z),T(w),T(z),T(q.$$.fragment,Z),T(U.$$.fragment,Z),se=!0)},o(Z){v(C),v(i),v(F),v(c.$$.fragment,Z),v(w),v(z),v(q.$$.fragment,Z),v(U.$$.fragment,Z),se=!1},d(Z){Z&&I(t),C&&C.d(),Y[r].d(),F&&F.d(),X(c),~m&&M[m].d(),z&&z.d(),X(q),X(U),ue=!1,ut(ye)}}}function Pt(o){let e,t;return e=new Tn({props:{open:o[4]}}),{c(){G(e.$$.fragment)},m(n,l){J(e,n,l),t=!0},p(n,l){const s={};l[0]&16&&(s.open=n[4]),e.$set(s)},i(n){t||(T(e.$$.fragment,n),t=!0)},o(n){v(e.$$.fragment,n),t=!1},d(n){X(e,n)}}}function jt(o){let e,t;return e=new zn({props:{open:o[2].open,type:o[2].type,msg:o[2].msg,timer:2e4}}),e.$on("dismiss",o[52]),{c(){G(e.$$.fragment)},m(n,l){J(e,n,l),t=!0},p(n,l){const s={};l[0]&4&&(s.open=n[2].open),l[0]&4&&(s.type=n[2].type),l[0]&4&&(s.msg=n[2].msg),e.$set(s)},i(n){t||(T(e.$$.fragment,n),t=!0)},o(n){v(e.$$.fragment,n),t=!1},d(n){X(e,n)}}}function Zt(o){let e,t;const n=[o[3]];let l={};for(let s=0;s<n.length;s+=1)l=En(l,n[s]);return e=new vn({props:l}),e.$on("dismiss",o[53]),{c(){G(e.$$.fragment)},m(s,r){J(e,s,r),t=!0},p(s,r){const i=r[0]&8?Rn(n,[Cn(s[3])]):{};e.$set(i)},i(s){t||(T(e.$$.fragment,s),t=!0)},o(s){v(e.$$.fragment,s),t=!1},d(s){X(e,s)}}}function Cl(o){let e,t,n,l,s,r,i,a,f,u,d,c,_,h,m,w,R,A,x,j,te,q=[],W=new Map,xe,ne,le,U,me,be,se,ue,ye,C;const Q=[dl,pl],Y=[];function Te(E,z){return E[6]?0:1}l=Te(o),s=Y[l]=Q[l](o),f=new ae({props:{type:"outline",icon:"plus",path:"icons"}}),h=new ae({props:{type:"outline",icon:"plus",path:"icons"}}),x=new ae({props:{type:"outline",icon:"save",path:"icons"}});let Re=$e(o[5].cells);const F=E=>E[67].id;for(let E=0;E<Re.length;E+=1){let z=vt(o,Re,E),fe=F(z);W.set(fe,q[E]=Bt(fe,z))}let V=o[4]&&Pt(o),H=o[2].open&&jt(o),M=o[3].open&&Zt(o);return{c(){e=y("div"),t=y("div"),n=y("button"),s.c(),r=D(),i=y("button"),a=y("span"),G(f.$$.fragment),u=ge(" Code"),d=D(),c=y("button"),_=y("span"),G(h.$$.fragment),m=ge(" Markdown"),w=D(),R=y("button"),A=y("span"),G(x.$$.fragment),j=D(),te=y("div");for(let E=0;E<q.length;E+=1)q[E].c();xe=D(),ne=y("style"),ne.textContent=`:global(body) {
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
    }`,le=D(),V&&V.c(),U=D(),H&&H.c(),me=D(),M&&M.c(),be=D(),se=y("div"),g(n,"class","btn btn-ghost btn-xs"),n.disabled=o[6],g(n,"title","Execute All Cells"),he(n,"executing",o[6]),g(a,"class","w-5 h-5"),g(i,"class","btn btn-ghost btn-xs"),g(i,"title","Add SQL Cell"),g(_,"class","w-5 h-5"),g(c,"class","btn btn-ghost btn-xs"),g(c,"title","Add Markdown Cell"),g(A,"class","w-5 h-5"),g(R,"class","btn btn-ghost btn-xs"),g(R,"title","Save Notebook"),g(t,"class","notebook-toolbar mt-2"),g(te,"class","notebook-content"),g(e,"class","--flex --flex-col"),g(se,"class","input-primary input-lg"),We(se,"display","none")},m(E,z){$(E,e,z),k(e,t),k(t,n),Y[l].m(n,null),k(t,r),k(t,i),k(i,a),J(f,a,null),k(i,u),k(t,d),k(t,c),k(c,_),J(h,_,null),k(c,m),k(t,w),k(t,R),k(R,A),J(x,A,null),k(e,j),k(e,te);for(let fe=0;fe<q.length;fe+=1)q[fe]&&q[fe].m(te,null);k(e,xe),k(e,ne),$(E,le,z),V&&V.m(E,z),$(E,U,z),H&&H.m(E,z),$(E,me,z),M&&M.m(E,z),$(E,be,z),$(E,se,z),ue=!0,ye||(C=[ee(n,"click",o[20]),ee(i,"click",o[32]),ee(c,"click",o[33]),ee(R,"click",o[34])],ye=!0)},p(E,z){let fe=l;l=Te(E),l!==fe&&(oe(),v(Y[fe],1,1,()=>{Y[fe]=null}),ie(),s=Y[l],s||(s=Y[l]=Q[l](E),s.c()),T(s,1),s.m(n,null)),(!ue||z[0]&64)&&(n.disabled=E[6]),(!ue||z[0]&64)&&he(n,"executing",E[6]),z[0]&267124128&&(Re=$e(E[5].cells),oe(),q=_n(q,z,F,1,E,Re,W,te,Sn,Bt,null,vt),ie()),E[4]?V?(V.p(E,z),z[0]&16&&T(V,1)):(V=Pt(E),V.c(),T(V,1),V.m(U.parentNode,U)):V&&(oe(),v(V,1,1,()=>{V=null}),ie()),E[2].open?H?(H.p(E,z),z[0]&4&&T(H,1)):(H=jt(E),H.c(),T(H,1),H.m(me.parentNode,me)):H&&(oe(),v(H,1,1,()=>{H=null}),ie()),E[3].open?M?(M.p(E,z),z[0]&8&&T(M,1)):(M=Zt(E),M.c(),T(M,1),M.m(be.parentNode,be)):M&&(oe(),v(M,1,1,()=>{M=null}),ie())},i(E){if(!ue){T(s),T(f.$$.fragment,E),T(h.$$.fragment,E),T(x.$$.fragment,E);for(let z=0;z<Re.length;z+=1)T(q[z]);T(V),T(H),T(M),ue=!0}},o(E){v(s),v(f.$$.fragment,E),v(h.$$.fragment,E),v(x.$$.fragment,E);for(let z=0;z<q.length;z+=1)v(q[z]);v(V),v(H),v(M),ue=!1},d(E){E&&(I(e),I(le),I(U),I(me),I(be),I(se)),Y[l].d(),X(f),X(h),X(x);for(let z=0;z<q.length;z+=1)q[z].d();V&&V.d(E),H&&H.d(E),M&&M.d(E),ye=!1,ut(C)}}}function Oe(){return`cell_${Date.now()}_${Math.random().toString(36).substr(2,9)}`}function Sl(o){const e=Object.keys(o[0]),t=[e.join(","),...o.map(r=>e.map(i=>`"${r[i]}"`).join(","))].join(`
`),n=new Blob([t],{type:"text/csv"}),l=URL.createObjectURL(n),s=document.createElement("a");s.href=l,s.download="notebook_data.csv",s.click(),URL.revokeObjectURL(l)}function $l(o,e,t){let n,l;_t(o,wn,p=>t(54,n=p)),_t(o,xn,p=>t(55,l=p));let{table:s}=e,{data:r={notebook:"Untitled",active:!0}}=e,{option:i={}}=e,{action:a={}}=e,f={open:!1,msg:null,type:null},u={open:!1,width:12,heigth:85,header:!0,title:"Dialog Title",header_actions:[],component:null},d=!1,c={name:"Sample Data Analysis Notebook",description:"A sample notebook demonstrating SQL queries and markdown documentation",useDuckDB:!0,cells:[{id:Oe(),isExecuted:!1,type:"markdown",code:`# Data Analysis Notebook

This notebook demonstrates how to query and analyze data using SQL.

## Getting Started

You can write **markdown** with *formatting*, lists, and more:

- Execute cells individually
- Run all cells at once
- Export results as CSV
- Save your notebook for later
- Execute them outomatic on a cron job defined on the Admin/Cron Jobs to nb/name_of_the_notebook

Powered by DuckDB`,showCode:!0,showOutput:!0,isEditing:!1},{id:Oe(),isExecuted:!1,type:"markdown",code:`# CONNECTING TO A DATA SOURCE
        Here for example we are connecting to the data lake for the exemple on etlx`,showCode:!0,showOutput:!0,isEditing:!1},{id:Oe(),isExecuted:!1,type:"sql",code:`INSTALL sqlite;
LOAD sqlite;
ATTACH IF NOT EXISTS 'ducklake:sqlite:database/dl_metadata.sqlite' AS dl (DATA_PATH 'database/dl');`,showCode:!0,showOutput:!0,isEditing:!1},{id:Oe(),isExecuted:!1,type:"sql",code:"SELECT * FROM dl.lineitem LIMIT 10;",showCode:!0,showOutput:!0,isEditing:!1},{id:Oe(),isExecuted:!1,type:"sql",code:"DETACH dl;",showCode:!0,showOutput:!0,isEditing:!1}]},_=!1,h=null,m=null;function w(p){const b={id:Oe(),isExecuted:!1,type:p,code:"",showCode:!0,showOutput:!0,isEditing:p==="sql"};t(5,c.cells=[...c.cells,b],c)}function R(p,b){const S={id:Oe(),isExecuted:!1,type:p,code:"",showCode:!0,showOutput:!0,isEditing:p==="sql"};t(5,c.cells=[...c.cells.slice(0,b+1),S,...c.cells.slice(b+1)],c)}function A(p){t(5,c.cells=c.cells.filter(b=>b.id!==p),c)}function x(p,b){const S=c.cells.findIndex(O=>O.id===p);S!==-1&&t(5,c.cells[S].code=b,c)}function j(p){const b=c.cells.findIndex(S=>S.id===p);b!==-1&&t(5,c.cells[b].isEditing=!c.cells[b].isEditing,c)}function te(p){const b=c.cells.findIndex(S=>S.id===p);b!==-1&&t(5,c.cells[b].isEditing=!0,c)}function q(p){try{return N(p)}catch(b){return`<p style="color: #ef4444;">Error rendering markdown: ${b.message}</p>`}}function W(p){const b=c.cells.findIndex(S=>S.id===p);b!==-1&&t(5,c.cells[b].showCode=!c.cells[b].showCode,c)}function xe(p){const b=c.cells.findIndex(S=>S.id===p);b!==-1&&t(5,c.cells[b].showOutput=!c.cells[b].showOutput,c)}const ne=async p=>{var O,ce,K,pe,Ee,Ce,ve,ze,Ae,de,Qe,Xe,De,Ye;const b=performance.now(),S={class:"nb",method:"cells",conf:n==null?void 0:n.conf,token:n==null?void 0:n.token,app:n==null?void 0:n.selected_app,data:{...p}};t(4,d=!0);try{const L=await In(S);if(L.success===!0)for(let Ne=0;Ne<c.cells.length;Ne++){const Le=(O=c.cells[Ne])==null?void 0:O.id;if((ce=L==null?void 0:L.data)!=null&&ce[Le]){const fn=Be((pe=(K=L==null?void 0:L.data)==null?void 0:K[Le])==null?void 0:pe.start).clone(),pn=Be((Ce=(Ee=L==null?void 0:L.data)==null?void 0:Ee[Le])==null?void 0:Ce.end).clone(),dn=Be.duration(pn.diff(fn)).clone();t(5,c.cells[Ne].isExecuted=!0,c),t(5,c.cells[Ne].isExecuting=!1,c),t(5,c.cells[Ne].output={data:Se.cloneDeep((ze=(ve=L==null?void 0:L.data)==null?void 0:ve[Le])==null?void 0:ze.data),error:((de=(Ae=L==null?void 0:L.data)==null?void 0:Ae[Le])==null?void 0:de.success)===!0?null:(Xe=(Qe=L==null?void 0:L.data)==null?void 0:Qe[Le])==null?void 0:Xe.msg,executionTime:(Ye=(De=L==null?void 0:L.data)==null?void 0:De[Le])!=null&&Ye.end?dn:Math.round(performance.now()-b)},c)}}else t(2,f.open=!0,f),t(2,f.type="error",f),t(2,f.msg=L.msg||L.toString(),f)}catch(L){console.log(L,L==null?void 0:L.message),t(2,f.open=!0,f),t(2,f.type="error",f),t(2,f.msg=`${L==null?void 0:L.message}`,f)}t(4,d=!1)},le=p=>p.map(b=>(delete b.output,b));async function U(p){const b=c.cells.findIndex(K=>K.id===p);if(b===-1)return;const S=c.cells[b];S.isExecuting=!0,t(5,c.cells=[...c.cells],c);const O=Se.cloneDeep(c.cells).filter((K,pe)=>pe<=b&&(K==null?void 0:K.isExecuted)||pe===b),ce=performance.now();try{await new Promise(K=>setTimeout(K,Math.random()*1e3+500)),S.type==="sql"?await ne({cells:le(O)}):S.output={executionTime:Math.round(performance.now()-ce)},S.isExecuted=!0}catch(K){S.output={error:K.message,executionTime:Math.round(performance.now()-ce)}}finally{S.isExecuting=!1,t(5,c.cells=[...c.cells],c)}}async function me(){t(6,_=!0),await ne({...Se.cloneDeep(c),cells:le(Se.cloneDeep(c==null?void 0:c.cells))}),t(6,_=!1)}async function be(p){const b=c.cells.findIndex(O=>O.id===p),S=[];for(let O=0;O<b;O++)S.push(Se.cloneDeep(c.cells[O]));await ne({cells:le(S)})}async function se(p){const b=c.cells.findIndex(O=>O.id===p),S=[];for(let O=b+1;O<c.cells.length;O++)S.push(Se.cloneDeep(c.cells[O]));await ne({cells:le(S)})}function ue(p,b){t(7,h=b),p.dataTransfer.effectAllowed="move"}function ye(p,b){p.preventDefault(),p.dataTransfer.dropEffect="move",t(8,m=b)}function C(){t(8,m=null)}function Q(p,b){if(p.preventDefault(),h!==null&&h!==b){const S=c.cells[h],O=[...c.cells];O.splice(h,1);const ce=h<b?b-1:b;O.splice(ce,0,S),t(5,c.cells=O,c)}t(7,h=null),t(8,m=null)}function Y(){t(7,h=null),t(8,m=null)}const Te=async p=>{var O,ce,K,pe,Ee,Ce;p.no_loding||t(4,d=!0);const b={lang:n==null?void 0:n.lang,conf:n==null?void 0:n.conf,token:n==null?void 0:n.token,app:n==null?void 0:n.selected_app,table:tables,filters:[{field:"etl_report_base_id",cond:"=",value:(p==null?void 0:p.etl_report_base_id)||(c==null?void 0:c.etl_report_base_id)},{field:"ref",cond:"LIKE",value:Be((O=steps[0])==null?void 0:O.date_ref).format("YYYY-MM-DD")}],ignore_filter:{etl_report_base:["etl_report_base_id","database"]},limit:-1},S=await An(b);if(console.log("data:",S),S.success===!0&&S.data){(pe=(K=(ce=S.data)==null?void 0:ce[tables[1]])==null?void 0:K.data)==null||pe.forEach(ve=>{var ze;ve.date_ref=(ze=steps==null?void 0:steps[0])==null?void 0:ze.date_ref}),steps[3].show_html=!1;for(const ve in S.data)Object.hasOwnProperty.call(S.data,ve)&&((Ce=(Ee=S.data[ve])==null?void 0:Ee.data)==null||Ce.map((ze,Ae)=>(ze.index=Ae,ze)));db_data={...S.data}}else t(2,f.open=!0,f),t(2,f.type="error",f),t(2,f.msg=S.msg,f);t(4,d=!1)},Re=p=>{console.log("after_save_form".toUpperCase(),p),Te({})},F=p=>{var b,S,O,ce,K,pe;if((O=(S=(b=n==null?void 0:n.tables)==null?void 0:b[p])==null?void 0:S.custom_form)!=null&&O.config)try{return JSON.parse((pe=(K=(ce=n==null?void 0:n.tables)==null?void 0:ce[p])==null?void 0:K.custom_form)==null?void 0:pe.config)}catch(Ee){return console.log(Ee.message),{}}},V=async p=>{var b,S,O,ce,K,pe,Ee,Ce,ve,ze,Ae;try{const de=Se.cloneDeep(p.data||{});p.clone===!0&&(de._original_pk=de[(S=(b=n==null?void 0:n.tables)==null?void 0:b[p.table])==null?void 0:S.pk],delete de[(ce=(O=n==null?void 0:n.tables)==null?void 0:O[p.table])==null?void 0:ce.pk],de.clone=!0);let Qe=`${l.t("crud.add")} - ${(pe=(K=n==null?void 0:n.tables)==null?void 0:K[p.table])==null?void 0:pe.comment}`;de!=null&&de[(Ce=(Ee=n==null?void 0:n.tables)==null?void 0:Ee[p.table])==null?void 0:Ce.pk]&&(Qe=`${l.t("crud.edit")} - ${(ze=(ve=n==null?void 0:n.tables)==null?void 0:ve[p.table])==null?void 0:ze.comment}`),t(4,d=!0);const Xe=(await $n(async()=>{const{default:Ye}=await import("./GlobalForm-CfFOlcE5.js");return{default:Ye}},__vite__mapDeps([0,1,2,3,4]),import.meta.url)).default;t(3,u.actions=null,u),t(3,u.header=!0,u),t(3,u.component=Xe,u),t(3,u.msg=null,u),t(3,u.html_msg=null,u),t(3,u.header_actions=[{type:"icon",icon:"cog-8-tooth",name:"form_customization",label:l.t("crud.form_customization")},{type:"icon",icon:"save",name:"save",label:l.t("crud.save"),action:Re}],u),t(3,u.title=Qe,u),t(3,u.data={table:p==null?void 0:p.table,data:de,aux_data:Se.cloneDeep(p==null?void 0:p.aux_data),main_table:(p==null?void 0:p.main_table)||s},u);const De=F(p==null?void 0:p.table);t(3,u.width=((Ae=De==null?void 0:De.layout)==null?void 0:Ae.size)||10,u),t(3,u.heigth=null,u),t(3,u.open=!0,u),t(4,d=!1)}catch(de){t(4,d=!1),console.log(de.message),t(2,f.open=!0,f),t(2,f.type="error",f),t(2,f.msg=de.message,f)}},H=p=>{t(3,u.actions=null,u),t(3,u.open=!1,u)};yn(async()=>{l.changeLanguage(n==null?void 0:n.lang);try{Be.locale(n==null?void 0:n.lang)}catch{Be.locale("en-us")}console.log("NOTEBOOK  Mount:",{table:s,data:r,option:i});try{r&&r!=null&&r.notebook_conf&&t(5,c=JSON.parse(r==null?void 0:r.notebook_conf))}catch(p){console.log(p.message)}});const M=()=>w("sql"),E=()=>w("markdown"),z=()=>V({table:s,data:{...r||{notebook:"Untitled",active:!0},notebook_conf:JSON.stringify({...c,cells:le(c.cells)})}}),fe=p=>j(p.id),st=p=>W(p.id),rt=p=>U(p.id),ot=p=>xe(p.id),it=p=>be(p.id),Z=p=>se(p.id);function re(p,b){p[b].name=this.value,t(5,c)}const je=p=>A(p.id),Ze=(p,b)=>{x(p.id,b)},en=p=>te(p.id),tn=(p,b)=>{x(p.id,b)},nn=p=>Sl(p.output.data),ln=p=>R("sql",p),sn=p=>R("markdown",p),rn=(p,b)=>ue(b,p),on=(p,b)=>ye(b,p),an=(p,b)=>Q(b,p),cn=p=>{t(2,f.open=!1,f)},un=p=>H(p.detail);return o.$$set=p=>{"table"in p&&t(1,s=p.table),"data"in p&&t(0,r=p.data),"option"in p&&t(30,i=p.option),"action"in p&&t(31,a=p.action)},o.$$.update=()=>{o.$$.dirty[0]&2&&console.log({table:s}),o.$$.dirty[1]&1&&((a==null?void 0:a.name)=="CANCEL"&&(a==null||a.action()),(a==null?void 0:a.name)=="SAVE"&&console.log(a),(a==null?void 0:a.name)=="REFRESH"&&console.log(a),(a==null?void 0:a.name)=="RUN_ALL"&&console.log(a))},[r,s,f,u,d,c,_,h,m,w,R,A,x,j,te,q,W,xe,le,U,me,be,se,ue,ye,C,Q,Y,V,H,i,a,M,E,z,fe,st,rt,ot,it,Z,re,je,Ze,en,tn,nn,ln,sn,rn,on,an,cn,un]}class ql extends mn{constructor(e){super(),bn(this,e,$l,Cl,kn,{table:1,data:0,option:30,action:31},null,[-1,-1,-1])}}export{ql as default};
