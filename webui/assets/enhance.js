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

disableKeydownOnInputs();

function disableKeydownOnInputs() {
    document.querySelectorAll('input').forEach(element => {
	element.addEventListener('keydown', e => {
	    e.stopPropagation();
	    markDirty(element);   
	});
    });

    document.querySelectorAll('textarea').forEach(element => {
	element.addEventListener('keydown', e => {
	    e.stopPropagation();
	    markDirty(element);   
	});
    });
}

function markDirty(element) {
    if(formIsDirty(element.form)) {
	element.classList.add('dirty')
    } else {
	element.classList.remove('dirty')		
    }
}


// https://stackoverflow.com/questions/598951/what-is-the-easiest-way-to-detect-if-at-least-one-field-has-been-changed-on-an-h
/*
** Determines if a form is dirty by comparing the current
** value of each element with its default value.
**
** @param {Form} form the form to be checked.
** @return {Boolean} true if the form is dirty, false otherwise.
*/
function formIsDirty(form) {
    for (var i = 0; i < form.elements.length; i++) {
        var element = form.elements[i];
        var type = element.type;
        switch (element.type) {
        case "checkbox":
        case "radio":
            if (element.checked != element.defaultChecked)
                return true;
            break;
        case "number":
        case "hidden":
        case "password":
        case "date":
        case "text":
        case "textarea":
            if (element.value != element.defaultValue)
                return true;
            break;
        case "select-one":
        case "select-multiple":
            for (var j = 0; j < element.options.length; j++)
                if (element.options[j].selected != element.options[j].defaultSelected)
                    return true;
            break;
        }
    }
    return false;
}
function onBeforeUnload(event) {
        event = event || window.event;
        for (i = 0; i < document.forms.length; i++) {
            switch (document.forms[i].id) {
            case "search":
                break;
            default:
                if (formIsDirty(document.forms[i])) {
                    if (event)
                        event.returnValue = "You have unsaved changes.";
                    return "You have unsaved changes.";
                }
                break;
            }
        }
    }
