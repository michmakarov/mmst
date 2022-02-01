
"use strict";

//220110 10:30 
//The /home/mich412/go/src/mak_common/kjs/xhr.js have been taken as a prototype
//220111 04:10 The functions here (in general) are strongly bimded to the mmsite project


var mmsitRevision="221027 17:22"



//For what does it serve? I just have liked it, it's all
//Let it remains. Others I will remove and needed things will be worked out at scratch.
function makeid(len) {
  var text = "";
  var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

  for (var i = 0; i < len-1; i++)
    text += possible.charAt(Math.floor(Math.random() * possible.length));

  return text;
}

//210111 04:07
//This makes a get query with responseType="text" and setting the answer or error message as string into given elemment (which defined by the id, the parameter elid)
//220124 18:15 my be elid==undefined//220126 04:44 This is brilliant an example of absence of pondering or a flashy bad decision. So let's to play back.
function getQuery(uri, elid){
	var  el
	if (typeof elid !== 'undefined') {
		el=document.getElementById(elid);
	}else{
		throw "getQuery: no an element to receive the result"
	}
	//alert(el);
		var xhr = new XMLHttpRequest();
		xhr.responseType = "text";
		
		var onLoadFun = function(e){//success of waiting - response have come
			//console.log("execRequest :onLoadFun: e=="+e);
			if (xhr.status==200){
				el.innerHTML = xhr.responseText;
			}else{
					el.innerHTML="status=="+xhr.status+":"+xhr.responseText;
				};
		};
		var onErrorFun = function(e){
			//console.log("execRequest :onErrorFun: e=="+e);
			el.innerHTML = uri+" error:"+e;
		};

		xhr.onload = onLoadFun;
		xhr.onerror = onErrorFun;
		

		xhr.open("GET", uri);
		xhr.send();
		return xhr;
}

//220126 05:14 It performs the func if it obtained a status 200.
function getQueryFunc(uri, func){
	if (typeof func == 'undefined') {
		throw "getQueryFunc: no a function to work out the result"
	}
		var xhr = new XMLHttpRequest();
		//xhr.responseType = "text";
		
		var onLoadFun = function(e){//success of waiting - response have come
			if (xhr.status==200){
				func();
			}else{
					throw "getQueryFunc no 200; status=="+xhr.status+":"+xhr.responseText;
				};
		};
		var onErrorFun = function(e){
			throw "getQueryFunc error; url="+url+";err="+e;
		};

		xhr.onload = onLoadFun;
		xhr.onerror = onErrorFun;
		

		xhr.open("GET", uri);
		xhr.send();
		return xhr;
}

//220111 12:26
//220113 11:48
//220124 18:18
//220125 In searching langErr220124
//220126 06:52 The comments was leaved for historic resons
function changeLang(){
	var btn = document.getElementById("langchanging");
	var btnCont=btn.innerHTML;
	//alert("changeLang at beg: btn.innerHTML="+btn.innerHTML);
	var toDo = function(){
		document.location.reload();
	}
	switch (btnCont) {
		case "ru":
			//getQuery("/changeLang?en","langchanging");
			getQueryFunc("/changeLang?en", toDo);
		break;
		case "en":
			//getQuery("/changeLang?ru","langchanging");
			getQueryFunc("/changeLang?ru", toDo);
		break;
	}
	//btn = document.getElementById("langchanging");
	//document.location.reload(); it is mere foolishness! 220124 17:39 Why?
	//alert("changeLang: before reload: btnCont="+btnCont);
	//document.location.reload(true);
	//return false;
	//btnCont=btn.innerHTML;
	//alert("changeLang: after reload and getting new btnCont: btnCont="+btnCont);
	//btn.onchange=document.location.reload();
}


//220125 07:55
function setMmsitRevision(){
	var el = document.getElementById("mmsitRevision");
	el.innerHTML="Front(nnsit.js) revision="+mmsitRevision;
}




