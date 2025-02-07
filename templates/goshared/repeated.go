package goshared

const repTpl = `
	{{ $f := .Field }}{{ $r := .Rules }}

	{{ if $r.GetCoordinates }}
		if len({{ accessor .}}) != 2 {
			return {{ err . "coordinates value must contain exactly 2 item(s)" }}
		}
		coordinates := {{ accessor .}}
		var lng, lat = coordinates[0], coordinates[1]
		if lng > 180 || lng < -180 {
			return {{ err . "lng value must contain between -180 and 180" }}
		}
		if lat > 90 || lng < -90 {
			return {{ err . "lat value must contain between -90 and 90" }}
		}
	{{ end }}

	{{ if $r.GetMinItems }}
		{{ if eq $r.GetMinItems $r.GetMaxItems }}
			if len({{ accessor . }}) != {{ $r.GetMinItems }} {
				return {{ err . "value must contain exactly " $r.GetMinItems " item(s)" }}
			}
		{{ else if $r.MaxItems }}
			if l := len({{ accessor . }}); l < {{ $r.GetMinItems }} || l > {{ $r.GetMaxItems }} {
			 	return {{ err . "value must contain between " $r.GetMinItems " and " $r.GetMaxItems " items, inclusive" }}
			}
		{{ else }}
			if len({{ accessor . }}) < {{ $r.GetMinItems }} {
				return {{ err . "value must contain at least " $r.GetMinItems " item(s)" }}
			}
		{{ end }}
	{{ else if $r.MaxItems }}
		if len({{ accessor . }}) > {{ $r.GetMaxItems }} {
			return {{ err . "value must contain no more than " $r.GetMaxItems " item(s)" }}
		}
	{{ end }}

	{{ if $r.GetUnique }}
		{{ lookup $f "Unique" }} := {{ if isBytes $f.Type.Element -}}
			make(map[string]struct{}, len({{ accessor . }}))
		{{ else -}}
			make(map[{{ (typ $f).Element }}]struct{}, len({{ accessor . }}))
		{{ end -}}
	{{ end }}

	{{ if or $r.GetUnique (ne (.Elem "" "").Typ "none") }}
		for idx, item := range {{ accessor . }} {
			_, _ = idx, item
			{{ if $r.GetUnique }}
				if _, exists := {{ lookup $f "Unique" }}[{{ if isBytes $f.Type.Element }}string(item){{ else }}item{{ end }}]; exists {
					return {{ errIdx . "idx" "repeated value must contain unique items" }}
				} else {
					{{ lookup $f "Unique" }}[{{ if isBytes $f.Type.Element }}string(item){{ else }}item{{ end }}] = struct{}{}
				}
			{{ end }}

			{{ render (.Elem "item" "idx") }}
		}
	{{ end }}
`
