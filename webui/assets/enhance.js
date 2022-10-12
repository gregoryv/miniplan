document.addEventListener('keydown', function(event) {
    //console.log(event.keyCode, event.key)
    switch (event.keyCode) {
    case 68: // d
	window.scrollTo(0, document.body.scrollHeight);
	break;
	
    case 69: // e
	document.getElementById("nav1").click();
	break;

    case 70: // f
	window.scrollTo(0, 0);
	break;

    case 73: // i
	// insert at top
	document.getElementById("insert1").click();		
	break;

    case 82: // r
	document.getElementById("nav2").click();	
	break;
	
    case 65: // a
	var element = document.getElementById("lastInsertButton");

	if(typeof(element) != 'undefined' && element != null){
            document.getElementById("lastInsertButton").click();
	} 
	break;
    }
});

function disableKeydownOnInputs() {
    document.querySelectorAll('input').forEach(element => {
	element.addEventListener('keydown', e => {
	    e.stopPropagation();
	});
    });

    document.querySelectorAll('textarea').forEach(element => {
	element.addEventListener('keydown', e => {
	    e.stopPropagation();
	});
    });
}
