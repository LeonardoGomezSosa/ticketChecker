$("document").ready(function () {
    $("#AceptarCodigo").click(function () {
        x =  $("#Entrada").val()
        console.log("Entrada: " + x);
        var request = $.ajax({
            url: "/Timer",
            method: "POST",
            async: false,
            data: { Entrada: x },
            dataType: "html",
        });
        request.done(function (data) {
            $("body").html(data);
        });
        request.fail(function (data) {
            $("body").html(data);
        });
        request.always(function () {
            alert("Fue y volvio")
        });

    });
    $("#AceptarCodigo").keypress(function(){
        
    });
});