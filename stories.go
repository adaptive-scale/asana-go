package asana

import (
	"fmt"
	"time"
)

// StoryBase contains the text of a story, as used when creating a new comment
type StoryBase struct {
	// Human-readable text for the story or comment. This will
	// not include the name of the creator. Can be edited only if the story is a comment.
	//
	// Note: This is not guaranteed to be stable for a given type of story.
	// For example, text for a reassignment may not always say “assigned to
	// …”. The API currently does not provide a structured way of inspecting
	// the meaning of a story.
	Text string `json:"text,omitempty"`

	// HTML formatted text for a comment. This will not include the name of the creator.
	// Can be edited only if the story is a comment.
	//
	// Note: This field is only returned if explicitly requested using the
	// opt_fields query parameter.
	HTMLText string `json:"html_text,omitempty"`

	// Whether the story should be pinned on the resource.
	IsPinned bool `json:"is_pinned,omitempty"`
}

// Story represents an activity associated with an object in the Asana
// system. Stories are generated by the system whenever users take actions
// such as creating or assigning tasks, or moving tasks between projects.
// Comments are also a form of user-generated story.
//
// Stories are a form of history in the system, and as such they are read-
// only. Once generated, it is not possible to modify a story.
type Story struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	StoryBase

	// Read-only. The time at which this object was created.
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// True if the object is hearted by the authorized user, false if not.
	Hearted bool `json:"hearted,omitempty"`

	// Read-only. Array of users who have hearted this object.
	Hearts []*User `json:"hearts,omitempty"`

	// Read-only. The number of users who have hearted this object.
	NumHearts int32 `json:"num_hearts,omitempty"`

	// The user who created the story.
	CreatedBy *User `json:"created_by,omitempty"`

	// Read-only. The object this story is associated with. Currently may only
	// be a task.
	Target *Task `json:"target,omitempty"`

	// Read-only. The component of the Asana product the user used to trigger
	// the story.
	Source string `json:"source,omitempty"`

	// Read-only. The type of story this is.
	Type string `json:"type,omitempty"`
}

// Stories lists all stories attached to a task
func (t *Task) Stories(client *Client, opts ...*Options) ([]*Story, *NextPage, error) {
	client.trace("Listing stories for %q", t.Name)

	var result []*Story

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/tasks/%s/stories", t.ID), nil, &result, opts...)
	return result, nextPage, err
}

// CreateComment adds a comment story to a task
func (t *Task) CreateComment(client *Client, story *StoryBase) (*Story, error) {
	client.info("Creating comment for task %q", t.Name)

	result := &Story{}

	err := client.post(fmt.Sprintf("/tasks/%s/stories", t.ID), nil, result)
	return result, err
}

// UpdateStory updates the story and returns the full record for the updated story.
// Only comment stories can have their text updated, and only comment stories and attachment stories can be pinned.
// Only one of text and html_text can be specified.
func (s *Story) UpdateStory(client *Client, story *StoryBase) (*Story, error) {
	client.info("Updating story %s", s.ID)

	result := &Story{}

	err := client.put(fmt.Sprintf("/stories/%s", s.ID), nil, result)
	return result, err
}
