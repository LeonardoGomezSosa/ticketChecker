$("document").ready(function () {
    $("#aceptar").click(function () {
        entrada = $("#Entrada").val();
        ticket = $("#Ticket").val();
        surtidor = $("#Surtidor").val();


        var request = $.ajax({
            url: "/Timer",
            method: "POST",
            async: false,
            data: { Entrada: entrada, Ticket: ticket, Surtidor: surtidor },
            dataType: "html",
        });
        request.done(function (data) {
            $("body").html(data);
        });
        request.fail(function (data) {
            $("body").html(data);
        });
        request.always(function () {
            console.log("Fue y volvio");
            console.log("Entrada: " + entrada);
            console.log("Ticket: " + ticket);
            console.log("Surtidor: " + surtidor);
        });
    });

});