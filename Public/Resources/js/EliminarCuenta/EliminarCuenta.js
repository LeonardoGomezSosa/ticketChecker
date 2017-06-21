$("document").ready(function () {

$("#Loading").hide()
    $("#EliminarCatalogo").click(function () {
        Eliminar("Eliminar Catalogo", 0);
    });

    $("#EliminarCuenta").click(function () {
        Eliminar("Eliminar Cuenta", 1);
    });

    $("#Confirmar").click(function () {
        operacion = $("#op").val();
        Confirmar(operacion);
    });

});


function Eliminar(titulo, operacion) {
    $("#tituloModal").html(titulo);
    $("#op").val(operacion);
    $("#confirm").modal("show");
}

function Confirmar(operacion) {
    console.log("Operacion: " + operacion);
    var request = $.ajax({
        url: "/Eliminar",
        method: "POST",
        async: false,
        data: { op: operacion },
        dataType: "html",
    });

    request.done(function (data) {
        $("body").html(data);
    });
    request.fail(function (data ) {
        $("body").html(data);
    });
    request.always(function () {
        $("#confirm").modal("hide");
    });
}
