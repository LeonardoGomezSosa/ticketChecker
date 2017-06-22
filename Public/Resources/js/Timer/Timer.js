$("document").ready(function () {
    $("#AceptarCodigo").click(function () {
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
            entrada = $("#Entrada").val();
            ticket = $("#Ticket").val();
            surtidor = $("#Surtidor").val();
            if (ticket !== "" && surtidor === "" || ticket === "" && surtidor !== "") {
                console.log("Solo uno de los datos ha sido fijado")
                alert("Solo uno de los datos ha sido fijado")
            }
        });
    });

});