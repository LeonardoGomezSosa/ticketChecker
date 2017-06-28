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

            timerOn = $("#TimerOn").val();
            Concluido = myString == $("#Concluido").val();
            
            if (Concluido == false) {
                setTimeout(CloseAlert, 3000);
            }
        });

    });

});

function ShowAlert() {
    $("#TimerDiv").alert();
}
function CloseAlert() {
    $("#TimerDiv").alert('close');
}