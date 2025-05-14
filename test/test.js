// Ini adalah file JavaScript untuk pengujian
// Komentar ini akan dihapus

function showAlert() {
    var message = "Ini adalah pesan dari JavaScript eksternal!"; // Pesan
    console.log( message ) ;     // Tampilkan di console dengan spasi ekstra
    // alert(message); // Kita comment out alert agar tidak mengganggu testing
}

showAlert();

/*
  Blok komentar
  multi-baris
  yang juga akan
  dihapus.
*/

var   unusedVariable   =   true  ; // Variabel ini tidak digunakan

(function() {
    var x = 10;
    var y = 20; // Spasi ekstra
    var z = x + y;
    console.log( "Hasil penjumlahan: " + z )  ;
})();

// Fungsi lain dengan parameter
function addNumbers( num1 , num2 ) {
    // Ini komentar di dalam fungsi
    return num1 + num2 ;
}

var sum = addNumbers( 5 , 10 ) ;
console.log( "5 + 10 = " + sum );
