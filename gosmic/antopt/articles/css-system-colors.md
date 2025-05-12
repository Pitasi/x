---
title: "CSS System colors"
date: "2025-03-10"
description: "System colors are a way for the OS, the browser, or the users to provide some colors to your CSS. This is a great opportunity for accessibility, as the colors will also reflect users preferences on high contrast colors."
categories:
  - "css"
  - "web"
published: true
---

System colors are a way for the OS, the browser, or the users to provide some
colors to your CSS.

This is a great opportunity for accessibility, as the colors will also reflect
users preferences on high contrast colors.

With the exception of `AccentColor` and `AccentColorText`, there is wide
browser support. <sup>[[caniuse](https://caniuse.com/?search=system-colors)]</sup>

I didn't see this feature in the wild, but I've seen using browser specific
vars such as `--webkit-link`, which is a pity.

The `<system-color>` keywords are defined as follows, conveniently extracted by
the [CSS Spec](https://drafts.csswg.org/css-color/#css-system-colors) and
rendered by your current browser:

<div style="display: flex; flex-direction: column; gap: 20px; margin-top: 15px">
  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: AccentColor; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>AccentColor</strong>
      <span>Background of accented user interface controls. (unsupported)</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: AccentColorText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>AccentColorText</strong>
      <span>Text of accented user interface controls. (unsupported)</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: ActiveText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>ActiveText</strong>
      <span>Text in active links. For light backgrounds, traditionally red.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: ButtonBorder; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>ButtonBorder</strong>
      <span>The base border color for push buttons.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: ButtonFace; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>ButtonFace</strong>
      <span>The face background color for push buttons.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: ButtonText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>ButtonText</strong>
      <span>Text on push buttons.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: Canvas; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>Canvas</strong>
      <span>Background of application content or documents.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: CanvasText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>CanvasText</strong>
      <span>Text in application content or documents.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: Field; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>Field</strong>
      <span>Background of input fields.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: FieldText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>FieldText</strong>
      <span>Text in input fields.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: GrayText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>GrayText</strong>
      <span>Disabled text. (Often, but not necessarily, gray.)</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: Highlight; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>Highlight</strong>
      <span>Background of selected text, for example from ::selection.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: HighlightText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>HighlightText</strong>
      <span>Text of selected text.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: LinkText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>LinkText</strong>
      <span>Text in non-active, non-visited links. For light backgrounds, traditionally blue.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: Mark; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>Mark</strong>
      <span>Background of text that has been specially marked (such as by the HTML mark element).</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: MarkText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>MarkText</strong>
      <span>Text that has been specially marked (such as by the HTML mark element).</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: SelectedItem; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>SelectedItem</strong>
      <span>Background of selected items, for example a selected checkbox.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: SelectedItemText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>SelectedItemText</strong>
      <span>Text of selected items.</span>
    </div>
  </div>

  <div style="display: flex; flex-direction: row; align-items: center; gap: 10px">
    <div style="height:60px; aspect-ratio: 1; border-radius: 10px; background: VisitedText; border: 1px solid rgba(0,0,0,.1)"></div>
    <div style="display: flex; flex-direction: column;">
      <strong>VisitedText</strong>
      <span>Text in visited links. For light backgrounds, traditionally purple.</span>
    </div>
  </div>
</div>
