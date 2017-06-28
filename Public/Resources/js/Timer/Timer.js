$("document").ready(function () {
    
    $("#AceptarCodigo").on('click', function (e) {
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
            setTimeout(CloseAlert, 3000);
            if (ticket !== "" && surtidor === "" || ticket === "" && surtidor !== "") {
                console.log("Solo uno de los datos ha sido fijado")
            }
        });

    });

});

function ShowAlert(){
    $("#TimerDiv").alert();
}
function CloseAlert(){
    $("#TimerDiv").alert('close');
}