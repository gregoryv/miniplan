/* todo */
document.addEventListener('keydown', function(event) {
    console.log(event.key, event.keyCode);
    switch (event.keyCode) {
    case 82:
	window.scrollTo(0, 0); 
	break;
    case 65:
	document.getElementById("lastInsertButton").click();
	break;
    }
	
});
/*
r 82 
e 69 
s 83 
d 68 
f 70
*/
