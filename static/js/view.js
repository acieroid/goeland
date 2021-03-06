var TASK_HTML = '<div class="task">\
  <div class="input-append">\
    <input class="span6 name" type="text" value="" placeholder="Task name">\
    <input class="span1 name" type="text" value="0" placeholder="Estimated time">\
    <div class="btn-group">\
      <input type="hidden" class="status" value="Todo">\
      <button class="btn dropdown-toggle btn-status" data-toggle="dropdown">\
        Todo\
        <span class="caret"></span>\
      </button>\
      <ul class="dropdown-menu">\
        <li><a class="btn btn-danger change-status">Todo</a></li>\
        <li><a class="btn btn-warning change-status">Started</a></li>\
        <li><a class="btn btn-success change-status">Done</a></li>\
      </ul>\
    </div>\
  </div>\
  <a class="btn remove-task"><span class="icon-remove"></span></a>\
  <a class="btn add-sub-task"><span class="icon-plus"></span></a>\
  <a class="btn toggle-description"><span class="icon-chevron-down"></span></a></br>\
  <textarea rows="3" style="display:none;" class="description" placeholder="Task description"></textarea>\
  <ul class="unstyled tasks">\
  </ul>\
</div>';
function updateButtons() {
    var todo = 0;
    var started = 0;
    var done = 0;

    /* Add the corresponding class for the task statuses */
    $('.btn-status').map(function() {
        var content = $(this).html();
        var task = $(this).parent().parent().parent()
        $(this).removeClass('btn-danger');
        $(this).removeClass('btn-warning');
        $(this).removeClass('btn-success');
        if (content.match(/Todo/)) {
            $(this).addClass('btn-danger');
            $(task).css('border-color', '#da4f49')
            todo += 1;
        } else if (content.match(/Started/)) {
            $(this).addClass('btn-warning');
            $(task).css('border-color', '#faa732')
            started += 1;
        } else if (content.match(/Done/)) {
            $(this).addClass('btn-success');
            $(task).css('border-color', '#5bb75b')
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
    $('#tasks').append(TASK_HTML)
    updateButtons();
}

/* add a subtask to a task */
function addSubTask() {
    $(this).parent().find('.tasks').first().append(TASK_HTML)
    updateButtons();
}

/* remove a task from the task list */
function removeTask() {
    $(this).parent().remove();
}

/* find the first-level subtasks, given a selector */
function findTasks(selector) {
    return $(selector).find('.task').first().parent().children('.task').map(function (i, elem) {
        return {'Id': 0,
                'Name': $(elem).find('.name').val(),
                'Description': $(elem).find('.description').val(),
                'Status': $(elem).find('.status').val(),
                'EstimatedTime': new Number($(elem).find('.etime').val()),
                'Items': findTasks(elem)
               };
    }).get();
}

/* return the todo list encoded as a json string */
function getEncodedList() {
    return {'Id': $(location).attr('href').match(/\/view\/([a-zA-Z0-9]+)/)[1], /* TODO */
            'Name': $('#name').val(),
            'ModificationTime': new Number($('#mtime').val()),
            'EstimatedTime': new Number($('#etime').val()),
            'Items': findTasks(document),
            }
}

/* save the currently displayed todo list */
function save() {
    $('.message').html('<i class="icon-refresh icon-white"></i> Saving...');
    $('.message').removeClass('label-important');
    $('.message').removeClass('label-success');
    $('.message').removeClass('label-info');
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
               $('.message').delay(3000).fadeOut();
           });
}

/* toggle the display of the description of a task */
function toggleDescription() {
    $(this).parent().children('.description').toggle();
}

/* hide all the descriptions */
function hideAllDescriptions() {
    $('.description').hide();
}

/* Enable drag and drop */
function enableDragAndDrop() {
      $('#tasks').nestedSortable({
        items: '.task',
        listType: 'div',
        opacity: 0.6,
        cursor: 'move',
    });
    $('.message').html('<i class="icon-info-sign icon-white"></i> Drag and drop enabled, reload the page to disable');
    $('.message').removeClass('label-important');
    $('.message').removeClass('label-success');
    $('.message').removeClass('label-info');
    $('.message').addClass('label-info');
    $('.message').fadeIn(3000).delay(3000).fadeOut();
    $('#btn-drag').addClass('disabled');
}

$(document).ready(function() {
    updateButtons();

    $(document).on('click', '.change-status', updateStatus);
    $(document).on('click', '.remove-task', removeTask);
    $(document).on('click', '.add-sub-task', addSubTask);
    $(document).on('click', '.toggle-description', toggleDescription);
    $('#btn-add-task').click(addTask);
    $('#save').click(save);
    $('#btn-hide-descr').click(hideAllDescriptions);
    $('#btn-drag').click(enableDragAndDrop);
});
