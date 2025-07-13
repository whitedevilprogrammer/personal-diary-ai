export function getSelectedText() {
  const selection = window.getSelection();
  return selection.toString();
}

export function getSelectionRect() {
  const selection = window.getSelection();
  if (selection.rangeCount === 0) return null;
  const range = selection.getRangeAt(0).cloneRange();
  const rect = range.getBoundingClientRect();
  return rect;
}
