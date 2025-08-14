/*
 * Copyright (c) 2025, WSO2 LLC. (http://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package text

// TextLine represents a single line in the TextDocument.
type TextLine interface {
	LineNo() uint
	Text() string
	StartOffset() uint
	EndOffset() uint
	EndOffsetWithNewLines() uint
	Length() uint
	LengthWithNewLineChars() uint
}

type textLineImpl struct {
	lineNo               uint
	text                 string
	startOffset          uint
	endOffset            uint
	lengthOfNewLineChars uint
}

func NewTextLine(lineNo uint, text string, startOffset, endOffset, lengthOfNewLineChars uint) TextLine {
	return &textLineImpl{
		lineNo:               lineNo,
		text:                 text,
		startOffset:          startOffset,
		endOffset:            endOffset,
		lengthOfNewLineChars: lengthOfNewLineChars,
	}
}

func (tl textLineImpl) LineNo() uint {
	return tl.lineNo
}

func (tl textLineImpl) Text() string {
	return tl.text
}

func (tl textLineImpl) StartOffset() uint {
	return tl.startOffset
}

func (tl textLineImpl) EndOffset() uint {
	return tl.endOffset
}

func (tl textLineImpl) EndOffsetWithNewLines() uint {
	return tl.endOffset + tl.lengthOfNewLineChars
}

func (tl textLineImpl) Length() uint {
	return tl.endOffset - tl.startOffset
}

func (tl textLineImpl) LengthWithNewLineChars() uint {
	return tl.endOffset - tl.startOffset + tl.lengthOfNewLineChars
}
