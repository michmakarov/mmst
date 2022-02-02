
"use strict";

//22201 07:02 it is a descendant of freelancer/mmsite/mmsit.js 

//versionInfo = :::0.0-1-ge8a5e5d[branch_*main,commit_e8a5e5d]_220202_1607:::

//var mmsitRevision="22201 07:02"




//220201 07:05
//This makes a get query with responseType="text" and setting the answer or error message as string into given elemment (which defined by the id, the parameter elid)
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

//220201 07:12 It performs the func if it obtained a status 200.
//If it encounts an error it throws an exception
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

//220125 In searching langErr220124
//220126 06:52 The comments was leaved for historic resons
//220201 07:25
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


//220201 08:01
//220202 13:55
//function setMmsitRevision(){
//	var el = document.getElementById("mmsitRevision");
//	el.innerHTML="Front(nnsit.js) revision="+mmsitRevision;
//}




