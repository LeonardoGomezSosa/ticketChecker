window.addEventListener("beforeunload", function (event) {
    $('#Loading').show();
});

$( document ).ready(function() {
     $('#Loading').hide();
    console.log("hola");
});