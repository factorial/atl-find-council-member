export function checkResponse(response: Response): Response {
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response;
}

/* This helpful function taken from https://gist.github.com/MoOx/8614711 */
export function createElement(options) {
  var el, a, i;
  if (!options.tagName) {
    el = document.createDocumentFragment();
  } else {
    el = document.createElement(options.tagName);
    if (options.className) {
      el.className = options.className;
    }

    if (options.attributes) {
      for (a in options.attributes) {
        el.setAttribute(a, options.attributes[a]);
      }
    }

    if (options.html !== undefined) {
      el.innerHTML = options.html;
    }
  }

  if (options.text) {
    el.appendChild(document.createTextNode(options.text));
  }

  if (options.childs && options.childs.length) {
    for (i = 0; i < options.childs.length; i++) {
      el.appendChild(
        options.childs[i] instanceof window.HTMLElement
          ? options.childs[i]
          : createElement(options.childs[i])
      );
    }
  }

  return el;
}
