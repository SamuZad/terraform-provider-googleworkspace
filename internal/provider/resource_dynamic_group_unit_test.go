package googleworkspace

import "testing"

func TestNormalizeDynamicGroupQuery(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		query string
		want  string
	}{
		"spaced API wrapper": {
			query: "(user.organizations.exists(org, org.department == 'engineering')) && user.is_guest_user == false",
			want:  "user.organizations.exists(org, org.department == 'engineering')",
		},
		"unspaced API wrapper": {
			query: "(user.organizations.exists(org, org.department == 'engineering'))&&user.is_guest_user==false",
			want:  "user.organizations.exists(org, org.department == 'engineering')",
		},
		"mixed whitespace API wrapper": {
			query: " \t(user.organizations.exists(org, org.department == 'engineering')) \t&&  user.is_guest_user\t ==\n false  ",
			want:  "user.organizations.exists(org, org.department == 'engineering')",
		},
		"no suffix": {
			query: "user.organizations.exists(org, org.department == 'engineering')",
			want:  "user.organizations.exists(org, org.department == 'engineering')",
		},
		"user-authored parenthesized query": {
			query: "(user.organizations.exists(org, org.department == 'engineering'))",
			want:  "(user.organizations.exists(org, org.department == 'engineering'))",
		},
		"non-trailing guest user condition": {
			query: "(user.is_guest_user == false && user.organizations.exists(org, org.department == 'engineering'))",
			want:  "(user.is_guest_user == false && user.organizations.exists(org, org.department == 'engineering'))",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := normalizeDynamicGroupQuery(test.query); got != test.want {
				t.Errorf("normalizeDynamicGroupQuery() = %q, want %q", got, test.want)
			}
		})
	}
}
