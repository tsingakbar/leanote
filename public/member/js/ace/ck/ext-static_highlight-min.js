ace.define("ace/ext/static_highlight",["require","exports","module","ace/edit_session","ace/layer/text","ace/config","ace/lib/dom"],function(e,t,i){var n=e("../edit_session").EditSession,s=e("../layer/text").Text,a=".ace_static_highlight {font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', 'Droid Sans Mono', monospace;font-size: 12px;}.ace_static_highlight .ace_gutter {width: 25px !important;display: block;float: left;text-align: right;padding: 0 3px 0 0;margin-right: 3px;position: static !important;}.ace_static_highlight .ace_line { clear: both; }.ace_static_highlight .ace_gutter-cell {-moz-user-select: -moz-none;-khtml-user-select: none;-webkit-user-select: none;user-select: none;}.ace_static_highlight .ace_gutter-cell:before {content: counter(ace_line, decimal);counter-increment: ace_line;}.ace_static_highlight {counter-reset: ace_line;}",c=e("../config"),o=e("../lib/dom"),r=function(e,t,i){var n=e.className.match(/lang-(\w+)/),s=t.mode||n&&"ace/mode/"+n[1];if(!s)return!1;var a=t.theme||"ace/theme/textmate",c="",l=[];if(e.firstElementChild)for(var h=0,d=0;d<e.childNodes.length;d++){var g=e.childNodes[d];3==g.nodeType?(h+=g.data.length,c+=g.data):l.push(h,g)}else c=o.getInnerText(e),t.trim&&(c=c.trim());r.render(c,s,a,t.firstLineNumber,!t.showGutter,function(t){o.importCssString(t.css,"ace_highlight"),e.innerHTML=t.html;for(var n=e.firstChild.firstChild,s=0;s<l.length;s+=2){var a=t.session.doc.indexToPosition(l[s]),c=l[s+1],r=n.children[a.row];r&&r.appendChild(c)}i&&i()})};r.render=function(e,t,i,s,a,o){function l(){var n=r.renderSync(e,t,i,s,a);return o?o(n):n}var h=1,d=n.prototype.$modes;return"string"==typeof i&&(h++,c.loadModule(["theme",i],function(e){i=e,--h||l()})),"string"==typeof t&&(h++,c.loadModule(["mode",t],function(e){d[t]||(d[t]=new e.Mode),t=d[t],--h||l()})),--h||l()},r.renderSync=function(e,t,i,c,o){c=parseInt(c||1,10);var r=new n("");r.setUseWorker(!1),r.setMode(t);var l=new s(document.createElement("div"));l.setSession(r),l.config={characterWidth:10,lineHeight:20},r.setValue(e);for(var h=[],d=r.getLength(),g=0;d>g;g++)h.push("<div class='ace_line'>"),o||h.push("<span class='ace_gutter ace_gutter-cell' unselectable='on'></span>"),l.$renderLine(h,g,!0,!1),h.push("\n</div>");var u="<div class='"+i.cssClass+"'><div class='ace_static_highlight' style='counter-reset:ace_line "+(c-1)+"'>"+h.join("")+"</div></div>";return l.destroy(),{css:a+i.cssText,html:u,session:r}},i.exports=r,i.exports.highlight=r}),function(){ace.require(["ace/ext/static_highlight"],function(){})}();