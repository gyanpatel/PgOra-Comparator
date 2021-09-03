window.addEventListener('DOMContentLoaded', event => {
    // Simple-DataTables
    // https://github.com/fiduswriter/Simple-DataTables/wiki

    const dbcompresult = document.getElementById('dbcompresult');
    if (dbcompresult) {
        new simpleDatatables.DataTable(dbcompresult);
    }
});
