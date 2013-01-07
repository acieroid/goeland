function updateButtons() {
    var todo = 0;
    var started = 0;
    var done = 0;

    /* Add the corresponding class for the task statuses */
    $('.btn-status').map(function() {
        var content = $(this).html();
        $(this).removeClass('btn-danger');
        $(this).removeClass('btn-warning');
        $(this).removeClass('btn-success');
        if (content.match(/Todo/)) {
            $(this).addClass('btn-danger');
            todo += 1;
        } else if (content.match(/Started/)) {
            $(this).addClass('btn-warning');
            started += 1;
        } else if (content.match(/Done/)) {
            $(this).addClass('btn-success');
            done += 1;
        }
        $(this).css('color', '#000');
    });

    var total = todo + started + done;
    todo = (todo*100)/total;
    started = (started*100)/total;
    done = (done*100)/total;

    $('.progress').html('<div class="bar bar-success" style="width: ' + done + '%;"></div>\
    <div class="bar bar-warning" style="width: ' + started + '%;"></div>\
    <div class="bar bar-danger" style="width: ' + todo + '%;"></div>');
}

/* update a 'status' button */
function updateStatus() {
    var status = $(this).html() + ' <span class="caret"></span>';
    var button = $(this).parent().parent().parent().find('.btn-status');
    var input = $(this).parent().parent().parent().find('.status');
    button.html(status);
    input.val($(this).html());
    updateButtons();
}

/* add a task to the task list */
function addTask() {
    var n = new Number($('.task').length + 1)
    var task = '<div class="task">\
    <div class="input-append">\
      <input class="span4 name" type="text" value="" placeholder="Task name">\
      <div class="btn-group">\
        <input type="hidden" class="status" value="Todo">\
        <button class="btn dropdown-toggle btn-status" data-toggle="dropdown">\
          Todo\
          <span class="caret"></span>\
        </button>\
        <ul class="dropdown-menu">\
          <li><a href="#" class="btn btn-danger change-status">Todo</a></li>\
          <li><a href="#" class="btn btn-warning change-status">Started</a></li>\
          <li><a href="#" class="btn btn-success change-status">Done</a></li>\
        </ul>\
      </div>\
    </div>\
    <button class="btn remove-task" type="button">Remove task</button></br>\
    <textarea rows="3" class="description" placeholder="Task description"></textarea>\
    <br/><br/>\
  </div>';
    $('.tasks').append(task)
    updateButtons();
}

/* remove a task from the task list */
function removeTask() {
    $(this).parent().remove();
}

/* return the todo list encoded as a json string */
function getEncodedList() {
    return {'Id': $(location).attr('href').match(/\/view\/([a-zA-Z0-9]+)/)[1], /* TODO */
            'Name': $('#name').val(),
            'ModificationTime': new Number($('#mtime').val()),
            'Items': $('.task').map(function (i, elem) {
                return {'Name': $(elem).find('.name').val(),
                        'Description': $(elem).find('.description').val(),
                        'Status': $(elem).find('.status').val()};
            }).get()};
}

/* save the currently displayed todo list */
function save() {
    $('.message').html('<i class="icon-refresh icon-white"></i> Saving...');
    $('.message').removeClass('label-important');
    $('.message').removeClass('label-success');
    $('.message').addClass('label-info');
    $('.message').fadeIn(3000);
    $.post('/save',
           {
               'list': JSON.stringify(getEncodedList())
           }, function(data) {
               var icon
               if (data.match(/Error/)) {
                   $('.message').removeClass('label-info');
                   $('.message').addClass('label-important');
                   icon = '<i class="icon-remove icon-white"></i>'
               } else {
                   $('.message').removeClass('label-info');
                   $('.message').addClass('label-success');
                   icon = '<i class="icon-ok icon-white"></i>'
               }
               $('.message').html(icon + ' ' +data);
               $(".message").delay(3000).fadeOut();
           });
}

$(document).ready(function() {
    updateButtons();

    $('.change-status').live('click', updateStatus);
    $('.remove-task').live('click', removeTask);
    $('#btn-add-task').click(addTask);
    $('#save').click(save);

    $("#tasks").sortable({
        opacity: 0.6,
        cursor: 'move',
    });
});